# Configuration

In this example we will be assuming the binary is called `pomfcrypt` and the application is run from a TTY.

## CLI Arguments

```
usage: pomfcrypt [<flags>]

Flags:
  -h, --help                Show context-sensitive help (also try --help-long and --help-man).
  -v, --debug               Enable debug output
      --max-size=256000000  Set maximum file size in bytes
      --filename-length=4   Set random filename length
  -d, --directory=uploads   Upload directory
  -s, --salt="salt"         Set salt for encryption
```

### Debug
Debug enables verbose output of all operations.

### Maximum filesize
Sets the maximum filesize in bytes

### Filename length
For each upload, a unique file name is generated. This setting configures the length of the generated file name.

### Directory
Sets the directory where uploaded files land. This directory must exist.

### Salt
Sets the salt needed for encryption of uploaded files.