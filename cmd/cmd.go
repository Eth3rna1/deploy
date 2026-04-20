package cmd

import (
    "os"
    "io"
    "os/exec"
)

func Move(src string, dest string) error {
    return os.Rename(src, dest)
}

func Run(command string) error {
    cmd := exec.Command("cmd", "/C", command)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}

func CreateCopy(file string, newFile string) error {
    in, err := os.Open(file)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(newFile)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }

    return out.Sync()
}
