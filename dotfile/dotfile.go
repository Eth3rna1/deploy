package dotfile

import (
    "os"
    "fmt"
    "maps"
    "errors"
    "strings"
    // "deploy/cmd"
)

func IsDotFileExists(dotfile string) bool {
    _, err := os.Stat(dotfile)
    return err == nil
}

func CreateDotFile(dotfile string) error {
    file, err := os.Create(dotfile)

    if err != nil {
        return err
    }

    file.Close()

    return nil
}

func WriteDotFile(dotfile string, content string) error {
    err := os.WriteFile(dotfile, []byte(content), 0644)

    return err
}

func IsDotFileEmpty(dotfile string) bool {
    bytes, err := os.ReadFile(dotfile)
    if err != nil { return true }

    content := string(bytes)
    content = strings.TrimSpace(content)

    return len(content) == 0
}

func InjectEmptyVariables(dotfile string) error {
    return WriteDotFile(dotfile, `# Use this variable to import another deployfile
DEPLOYFILE_LOC=

# if 'DEPLOYFILE_LOC' is defined, then all defined variables
# within such deployfile will get imported and overwrite any
# previously defined variables if defined in such deployfile.
# Everything defined past this comment will get overwritten
# or defined.

# Location of the global binary directory
GLOBAL_BIN_DIR=

# Location of the base directory of the project
BASE_DIR=

# Location of the local binary directory relative to BASE_DIR
LOCAL_BIN_DIR=

# Location of the local scripts directory relative to BASE_DIR
SCRIPTS_DIR=

# The CMD command to compile the project
COMPILATION_CMD=

# Tells where the binary is located after compiling relative to BASE_DIR
BINARY_LOC=`)
}


func ExpandEnvVar(query string) string {
    query = os.ExpandEnv(query) // supporting UNIX syntax as well

    var buffer strings.Builder
    var env_buffer strings.Builder
    capture := false

    for _, char := range query {
        if char == '%' {
            capture = !capture
            env_buffer.WriteRune(char)

            if capture == false {
                // meaning it was previously true (we just closed a %VAR%)

                envStr := env_buffer.String()
                n := len(envStr)

                val, exists := os.LookupEnv(envStr[1:n-1])

                if exists {
                    buffer.WriteString(val)
                } else {
                    // not actually an env variable
                    buffer.WriteString(envStr)
                }

                env_buffer.Reset()
            }

            continue
        }

        if capture {
            env_buffer.WriteRune(char)
        } else {
            buffer.WriteRune(char)
        }
    }

    buffer.WriteString(env_buffer.String())

    return buffer.String()
}

func LoadDotFile(dotfile string) (map[string]string, error) {
    byte_content, err := os.ReadFile(dotfile)

    content := string(byte_content) // converts bytes to a string

    if err != nil { return nil, err }

    variables := make(map[string]string)

    lines := strings.Split(content, "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)

        if line == "" || !strings.Contains(line, "=") {
            continue
        }

        if line[0] == '#' { continue } // allowing for comments

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            error_msg := fmt.Sprintf("Error: \"%s\". Missing information")
            return nil, errors.New(error_msg)
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        value = ExpandEnvVar(value)

        if value == "" {
            fmt.Printf("WARNING: \"%s\" has not been defined\n", key)
            continue
        }

        // in case of different deployfile reference
        if key == "DEPLOYFILE_LOC" {
            // recursive action
            sub_variables, err := LoadDotFile(value)
            if err != nil { return sub_variables, err }

            maps.Copy(variables, sub_variables)
            continue
        }

        variables[key] = value
    }

    return variables, nil
}
