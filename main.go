package main

//--------------------------------------------------------------------------------------------------
// PyEnv
//
// Pyenv is a tool used to install and remove Python versions, and to create and destroy Python
// virtual environments.
//
// Pyenv is MacOS-only.
//
// The goal of this program is to produce a tarball with the pyenv executable in it.  We
// do that by creating a hidden installation of homebrew, and use homebrew to install
// pyenv, also in a way that is hidden to the rest of the laptop software.  Then we create a
// tarball of the hidden homebrew directory.
//
//--------------------------------------------------------------------------------------------------

import (
  "bufio"
  "fmt"
  "os"
  "os/exec"
  "os/user"
  "path"
  "path/filepath"
)

type Builder struct {
}

func (this Builder) build() {
  repoDir, _ :=  filepath.Abs(filepath.Dir(os.Args[0]))
  pyenvDir := path.Join(repoDir, "pyenv")
  homebrewDir := path.Join(pyenvDir, "homebrew")
  homebrewFilepath := path.Join(homebrewDir, "bin", "brew")
  currentUser, _ := user.Current();
  homeDir := currentUser.HomeDir    
  tarballDir := path.Join(homeDir, ".toolbox-tarballs")
  tarballFilepath := path.Join(tarballDir, "pyenv.tgz")

  // If the pyenv directory exists, remove it
  _, err := os.Stat(pyenvDir)
  if !os.IsNotExist(err) {
    os.RemoveAll(pyenvDir)
    if err != nil {
      fmt.Printf("Error removing pyenv directory. %s\n", err.Error())
      os.Exit(1)
    }
  }

  // Create the homebrew directory
  os.MkdirAll(homebrewDir, 0700)

  // Install homebrew using curl, and build using XCode
  url := "https://github.com/Homebrew/brew/tarball/master"
  curlCommand := fmt.Sprintf("curl -L %s | tar xz --strip 1 -C %s", url, homebrewDir)
  fmt.Printf("command: '%s'\n", curlCommand)
  bashCommand := exec.Command("bash", "-c", curlCommand)
  stderr, _ := bashCommand.StderrPipe()
  bashCommand.Start()
  
  scanner := bufio.NewScanner(stderr)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Printf("%s\n", line)
  }
  err = bashCommand.Wait()
  if err != nil {
    fmt.Printf("Error installing homebrew: %s\n", err.Error())
    os.Exit(1)
  }

  // Brew update
  homebrewCommand := fmt.Sprintf("%s update", homebrewFilepath)
  bashCommand = exec.Command("bash", "-c", homebrewCommand)
  stderr, _ = bashCommand.StderrPipe()
  bashCommand.Start()

  scanner = bufio.NewScanner(stderr)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Printf("%s\n", line)
  }
  err = bashCommand.Wait()
  if err != nil {
    fmt.Printf("Error running 'brew update'. %s\n", err.Error())
    os.Exit(1)
  }

  // Brew install pyenv
  homebrewCommand = fmt.Sprintf("%s install %s", homebrewFilepath, "pyenv")
  bashCommand = exec.Command("bash", "-c", homebrewCommand)
  stderr, _ = bashCommand.StderrPipe()
  bashCommand.Start()

  scanner = bufio.NewScanner(stderr)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Printf("%s\n", line)
  }
  err = bashCommand.Wait()
  if err != nil {
    fmt.Printf("Error running 'brew install'. %s\n", err.Error())
    os.Exit(1)
  }

  // Create the tarball directory if it does not exist
  _, err = os.Stat(tarballDir)
  if os.IsNotExist(err) {
    os.MkdirAll(tarballDir, 0700)
  }

  // If the tarball exists, remove it
  _, err = os.Stat(tarballFilepath)
  if !os.IsNotExist(err) {
    os.RemoveAll(tarballFilepath)
    if err != nil {
      fmt.Printf("Error removing tarball. %s\n", err.Error())
      os.Exit(1)
    }
  }

  // Tar the pgadmin directory
  parentDir := filepath.Dir(repoDir)
  command := fmt.Sprintf("tar -C  %s -czf %s %s", parentDir, tarballFilepath, repoDir)
  fmt.Printf("command: '%s'\n", command)
  _, err = exec.Command("bash", "-c", command).Output()
  if err != nil {
    fmt.Printf("error tarring pgadmin. %s\n", err.Error())
    os.Exit(1)
  }
}

func main() {
  var builder Builder
  builder.build()
}
