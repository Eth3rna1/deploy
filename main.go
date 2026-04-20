package main

import (
    "os"
    "fmt"
    "path/filepath"
    "deploy/cmd"
    "deploy/dotfile"
)

const DEPLOYFILE string = "./.deployfile"

func initializeDotFile(file string) (map[string]string, error) {
    if !dotfile.IsDotFileExists(file) {
        if err := dotfile.CreateDotFile(file); err != nil {
            return nil, err
        }

        if err := dotfile.InjectEmptyVariables(file); err != nil {
            return nil, err
        }
    }

    if dotfile.IsDotFileEmpty(file) {
        if err := dotfile.InjectEmptyVariables(file); err != nil {
            return nil, err
        }
    }
    
    variables, err := dotfile.LoadDotFile(file)

    if (err != nil) { return nil, err }

    return variables, nil
}

// Checks if the given directory is the
// current active working directory
func IsCWD(base_dir string) bool {
    cwd, err := os.Getwd()

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    absolute_base_dir, err:= filepath.Abs(base_dir)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return cwd == absolute_base_dir
}

func main() {
    args := os.Args

    if len(args) < 2 {
        fmt.Println("Please provide a <SRC DIR / SRC FILE / BINARY>")
        return
    }

    source := args[1]

    vars, err := initializeDotFile(DEPLOYFILE)

    if err != nil {
        fmt.Println(err)
        return
    }
    
    cwd, err := os.Getwd()

    if err != nil {
        fmt.Println(err)
        return
    }

    os.Chdir(vars["BASE_DIR"])

    // attempting to compile to a binary
    if err := cmd.Run(vars["COMPILATION_CMD"]); err != nil {
        fmt.Println(err)
        fmt.Println("Could not compile source code")
        return
    }

    fmt.Println("Compiled binary")

    os.Chdir(cwd)

    binary_base_name := filepath.Base(vars["BINARY_LOC"])

    // optional variable
    if global_bin_dir, ok := vars["GLOBAL_BIN_DIR"]; ok {
        // attempting to create a copy of the located binary to the global env path variable
        new_binary_dest := filepath.Join(global_bin_dir, binary_base_name)
        if err := cmd.CreateCopy(vars["BINARY_LOC"], new_binary_dest); err != nil {
            fmt.Println("Could not find binary")
            fmt.Println("Reason: ", err)
            return
        }

        fmt.Println("Created copy to global binary directory")
    }

    // optional variable
    if local_bin_dir, ok := vars["LOCAL_BIN_DIR"]; ok {
        // attempting to move to binary to its local bin directory
        new_binary_dest := filepath.Join(local_bin_dir, binary_base_name)
        if err := cmd.Move(vars["BINARY_LOC"], new_binary_dest); err != nil {
            fmt.Println("Could not move binary")
            fmt.Println("Reason: ", err)
            return
        }

        fmt.Println("Moved binary to local binary directory")
    }

    // optional variable
    if scripts_dir, ok := vars["SCRIPTS_DIR"]; ok {
        // attempting to move the source entry into the defined SCRIPTS directory
        source_base_name := filepath.Base(source)
        new_scripts_dest := filepath.Join(scripts_dir, source_base_name)
        if err := cmd.Move(source, new_scripts_dest); err != nil {
            fmt.Println("Could not move source to scripts directory")
            fmt.Println("Reason: ", err)
        }

        fmt.Println("Moved src to local source directory")
    }
}
