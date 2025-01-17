package util

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func FirstError(errs ...error) (err error) {
	for _, e := range errs {
		if err == nil {
			err = e
		} else if e != nil {
			wwlog.ErrorExc(e, "Unhandled error")
		}
	}

	return
}

func DirModTime(path string) (time.Time, error) {

	var lastTime time.Time
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fi, err := os.Stat(path)
		if err != nil {
			return nil
		}
		stat := fi.Sys().(*syscall.Stat_t)
		cur := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
		if cur.After(lastTime) {
			lastTime = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
		}

		return nil
	})
	if err != nil {
		return time.Time{}, err
	}
	return lastTime, nil
}

func PathIsNewer(source string, compare string) bool {
	wwlog.Debug("Comparing times on paths: '%s' - '%s'", source, compare)

	time1, err := DirModTime(source)
	if err != nil {
		wwlog.DebugExc(err, "")
		return false
	}

	time2, err := DirModTime(compare)
	if err != nil {
		wwlog.DebugExc(err, "")
		return false
	}

	return time1.Before(time2)
}

/*
Checks if given string is in slice. I yes returns true, false otherwise.
*/
func InSlice(slice []string, match string) bool {
	for _, val := range slice {
		if val == match {
			return true
		}
	}
	return false
}

func IsDir(path string) bool {
	wwlog.Debug("Checking if path exists as a directory: %s", path)

	if path == "" {
		return false
	}
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

func IsFile(path string) bool {
	wwlog.Debug("Checking if path exists as a file: %s", path)

	if path == "" {
		return false
	}

	if stat, err := os.Lstat(path); err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func ReadFile(path string) ([]string, error) {
	lines := []string{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	f.Close()
	return lines, nil
}

func ValidString(pattern string, expr string) bool {
	if b, _ := regexp.MatchString(expr, pattern); b {
		return true
	}
	return false
}

func FindFiles(path string) []string {
	var ret []string

	wwlog.Debug("Finding files in path: %s", path)

	err := filepath.Walk(path, func(location string, info os.FileInfo, err error) error {
		if err != nil {
			wwlog.Warn("Error walking path %s: %v", location, err)
			return err
		}

		// Get the relative path from the base directory
		relPath, relErr := filepath.Rel(path, location)
		if relErr != nil {
			wwlog.Warn("Error computing relative path for %s: %v", location, relErr)
			return relErr
		}

		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			wwlog.Debug("FindFiles() found directory: %s", relPath)
			ret = append(ret, relPath+"/")
		} else {
			wwlog.Debug("FindFiles() found file: %s", relPath)
			ret = append(ret, relPath)
		}

		return nil
	})
	if err != nil {
		wwlog.Warn("Error during file walk: %v", err)
		return ret
	}

	return ret
}

/*
Finds all files under a given directory with tar like include and ignore patterns.
/foo/*
will match /foo/baar/ and /foo/baar/sibling
*/
func FindFilterFiles(
	path string,
	includePattern []string,
	ignorePattern []string,
	ignore_xdev bool) (ofiles []string, err error) {
	wwlog.Debug("Finding files: %s include: %s ignore: %s", path, includePattern, ignorePattern)

	// Preprocess patterns to remove leading (and trailing) /, as we are handling relative paths
	for i, pattern := range ignorePattern {
		ignorePattern[i] = strings.Trim(pattern, "/")
	}

	// Convert the base path to an absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ofiles, fmt.Errorf("failed to resolve absolute path: %s: %w", path, err)
	}

	// Expand the include list
	var globedInclude []string
	for _, include := range includePattern {
		globed, err := filepath.Glob(filepath.Join(absPath, include))
		if err != nil {
			return ofiles, fmt.Errorf("failed to glob pattern %s: %w", include, err)
		}
		globedInclude = append(globedInclude, globed...)
	}

	var dev uint64
	if ignore_xdev {
		wwlog.Debug("Ignoring cross-device (xdev) files")
		pathStat, err := os.Stat(absPath)
		if err != nil {
			return ofiles, fmt.Errorf("failed to stat base path: %s: %w", absPath, err)
		}
		dev = pathStat.Sys().(*syscall.Stat_t).Dev
	}

	for _, inc := range globedInclude {
		wwlog.Debug("Processing include pattern: %s", inc)
		stat, err := os.Lstat(inc)
		if err != nil {
			return ofiles, fmt.Errorf("failed to stat include: %s: %w", inc, err)
		}

		if stat.IsDir() {
			// Walk the directory
			err = filepath.WalkDir(inc, func(location string, info fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				relPath, relErr := filepath.Rel(absPath, location)
				if relErr != nil {
					wwlog.Warn("Error computing relative path for %s: %v", location, relErr)
					return relErr
				}

				if relPath == "." {
					return nil
				}

				fsInfo, err := info.Info()
				if err != nil {
					return err
				}

				if ignore_xdev && fsInfo.Sys().(*syscall.Stat_t).Dev != dev {
					wwlog.Debug("Ignored (cross-device): %s", location)
					return nil
				}

				for _, ignoredPattern := range ignorePattern {
					if ignored, _ := filepath.Match(ignoredPattern, relPath); ignored {
						wwlog.Debug("Ignored %s due to pattern %s", relPath, ignoredPattern)
						if info.IsDir() {
							return filepath.SkipDir
						}
						return nil
					}
				}

				ofiles = append(ofiles, relPath)
				return nil
			})
			if err != nil {
				return ofiles, fmt.Errorf("error walking directory %s: %w", inc, err)
			}
		} else {
			// Add the file directly
			relPath, relErr := filepath.Rel(absPath, inc)
			if relErr != nil {
				wwlog.Warn("Error computing relative path for %s: %v", inc, relErr)
				return ofiles, relErr
			}
			ofiles = append(ofiles, relPath)
		}
	}

	return ofiles, nil
}

// ******************************************************************************
func ExecInteractive(command string, a ...string) error {
	wwlog.Debug("ExecInteractive(%s, %s)", command, a)
	c := exec.Command(command, a...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	return err
}

func SystemdStart(systemdName string) error {
	startCmd := fmt.Sprintf("systemctl restart %s", systemdName)
	enableCmd := fmt.Sprintf("systemctl enable %s", systemdName)

	wwlog.Debug("Setting up Systemd service: %s", systemdName)
	err := ExecInteractive("/bin/sh", "-c", startCmd)
	if err != nil {
		return fmt.Errorf("failed to run start cmd: %w", err)
	}
	err = ExecInteractive("/bin/sh", "-c", enableCmd)
	if err != nil {
		return fmt.Errorf("failed to run enable cmd: %w", err)
	}

	return nil
}

func CopyUIDGID(source string, dest string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	// root is always good, if we failt to get UID/GID of a file
	var UID int = 0
	var GID int = 0
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		UID = int(stat.Uid)
		GID = int(stat.Gid)
	}
	wwlog.Debug("Chown %d:%d '%s'", UID, GID, dest)
	err = os.Chown(dest, UID, GID)
	return err
}

func IncrementIPv4(start net.IP, inc uint) net.IP {
	ipv4 := start.To4()
	v4_int := uint(ipv4[0])<<24 + uint(ipv4[1])<<16 + uint(ipv4[2])<<8 + uint(ipv4[3])
	v4_int += inc
	v4_o3 := byte(v4_int & 0xFF)
	v4_o2 := byte((v4_int >> 8) & 0xFF)
	v4_o1 := byte((v4_int >> 16) & 0xFF)
	v4_o0 := byte((v4_int >> 24) & 0xFF)
	ipv4_new := net.IPv4(v4_o0, v4_o1, v4_o2, v4_o3)
	return ipv4_new
}

/*
Appending the lines to the given file
*/
func AppendLines(fileName string, lines []string) error {
	wwlog.Verbose("appending %v lines to %s", len(lines), fileName)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can't open file: %s: %w", fileName, err)
	}
	defer file.Close()
	for _, line := range lines {
		wwlog.Debug("Appending '%s' to %s", line, fileName)
		if _, err := file.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return fmt.Errorf("can't write to file: %s: %w", fileName, err)
		}

	}
	return nil
}

/*
******************************************************************************

	Create an archive using cpio
*/
func CpioCreate(
	rootdir string,
	ifiles []string,
	ofile string,
	format string,
	cpio_args ...string) (err error) {

	args := []string{
		"--quiet",
		"--create",
		"--directory", rootdir,
		"--format", format,
		"--file=" + ofile}

	args = append(args, cpio_args...)

	proc := exec.Command("cpio", args...)

	stdin, err := proc.StdinPipe()
	if err != nil {
		return err
	}

	err_in := make(chan error, 1)
	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, strings.Join(ifiles, "\n"))
		err_in <- err
	}()

	out, err := proc.CombinedOutput()
	if err != nil {
		wwlog.Error("cpio failed: %s", out)
	} else if len(out) > 0 {
		wwlog.Debug("cpio: %s", out)
	}

	return FirstError(err, <-err_in)
}

/*
******************************************************************************

	Compress a file using gzip or pigz
*/
func FileGz(
	file string) (err error) {

	file_gz := file + ".gz"

	if IsFile(file_gz) {
		err := os.Remove(file_gz)

		if err != nil {
			return fmt.Errorf("could not remove existing file: %s: %w", file_gz, err)
		}
	}

	compressor, err := exec.LookPath("pigz")
	if err != nil {
		wwlog.Verbose("Could not locate PIGZ")
		compressor, err = exec.LookPath("gzip")
		if err != nil {
			wwlog.Verbose("Could not locate GZIP")
			return fmt.Errorf("no compressor program for image file: %s: %w", file_gz, err)
		}
	}

	wwlog.Verbose("Using compressor program: %s", compressor)

	proc := exec.Command(
		compressor,
		"--keep",
		file)

	out, err := proc.CombinedOutput()
	if len(out) > 0 {
		outStr := string(out[:])
		if err != nil && strings.HasSuffix(compressor, "gzip") && strings.Contains(outStr, "unrecognized option") {
			var gzippedFile *os.File
			var gzipStderr io.ReadCloser

			/* Older version of gzip, try it another way: */
			wwlog.Verbose("%s does not recognize the --keep flag, trying redirected stdout", compressor)

			/* Open the output file for writing: */
			gzippedFile, err = os.Create(file_gz)
			if err != nil {
				return fmt.Errorf("unable to open compressed image file for writing: %s: %w", file_gz, err)
			}

			/* We'll execute gzip with output to stdout and attach stdout to the compressed file we just
			   created:
			*/
			proc = exec.Command(
				compressor,
				"--stdout",
				file)
			proc.Stdout = gzippedFile
			gzipStderr, err = proc.StderrPipe()
			if err != nil {
				return fmt.Errorf("unable to open stderr pipe for compression program: %s: %w", compressor, err)
			}

			/* Execute the command: */
			err = proc.Start()
			if err != nil {
				_ = proc.Wait()
				gzippedFile.Close()
				os.Remove(file_gz)
				err = fmt.Errorf("unable to successfully execute compression program: %s: %w", compressor, err)
			} else {
				err = proc.Wait()
				gzippedFile.Close()
				if err != nil {
					os.Remove(file_gz)
					err = fmt.Errorf("unable to successfully create compressed image file: %s: %w", file_gz, err)
				} else {
					wwlog.Verbose("Successfully compressed image file: %s", file_gz)
				}
			}
			out, _ = io.ReadAll(gzipStderr)
		}
		wwlog.Debug(string(out))
	}

	return err
}

/*
******************************************************************************

	Create an archive using cpio
*/
func BuildFsImage(
	name string,
	rootfsPath string,
	imagePath string,
	include []string,
	ignore []string,
	ignore_xdev bool,
	format string,
	cpio_args ...string) (err error) {

	err = os.MkdirAll(path.Dir(imagePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create image directory for %s: %s: %w", name, imagePath, err)
	}
	wwlog.Debug("Created image directory for %s: %s", name, imagePath)

	files, err := FindFilterFiles(
		rootfsPath,
		include,
		ignore,
		ignore_xdev)
	if err != nil {
		return fmt.Errorf("failed discovering files for %s: %s: %w", name, rootfsPath, err)
	}

	err = CpioCreate(
		rootfsPath,
		files,
		imagePath,
		format,
		cpio_args...)
	if err != nil {
		return fmt.Errorf("failed creating image for %s: %s: %w", name, imagePath, err)
	}

	wwlog.Info("Created image for %s: %s", name, imagePath)

	err = FileGz(imagePath)
	if err != nil {
		return fmt.Errorf("failed to compress image for %s: %s: %w", name, imagePath+".gz", err)
	}

	wwlog.Info("Compressed image for %s: %s", name, imagePath+".gz")

	return nil
}

/*
Get size of given directory in bytes
*/
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

/*
Convert bytes to human friendly format
*/
func ByteToString(b int64) string {
	const base = 1024
	if b < base {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(base), 0
	for n := b / base; n >= base; n /= base {
		div *= base
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func HashFile(file *os.File) (string, error) {
	if prevOffset, err := file.Seek(0, 0); err != nil {
		return "", err
	} else {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
		if _, err := file.Seek(prevOffset, 0); err != nil {
			return "", err
		}
		return hex.EncodeToString(hasher.Sum(nil)), nil
	}
}

func EncodeYaml(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err := encoder.Encode(data)
	return buf.Bytes(), err
}

func EqualYaml(a interface{}, b interface{}) (bool, error) {
	aYaml, err := EncodeYaml(a)
	if err != nil {
		return false, err
	}

	bYaml, err := EncodeYaml(b)
	if err != nil {
		return false, err
	}

	return bytes.Equal(aYaml, bYaml), nil
}
