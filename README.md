# Upgraded-spork (ElvUI Updater)

This Go program downloads and extracts the latest version of the ElvUI addon for World of Warcraft from the official GitHub repository.

## Simple usage

You can download the binaries [here](<https://github.com/qqMelon/upgraded-spork/releases/tag/v1.0.0).
Or download the binaries you want on `bin/` forlder in project.

## development usage

1. Clone the repository:

  ```bash
  git clone https://github.com/your-username/elvui-updater.git
  ```

2. Navigate to the project directory:

  ```bash
  cd upgraded-spork
  ```

3. Run the program:

  ```bash
  go run main.go
  ```

The program will download the latest version of ElvUI, extract it, and place the relevant folders (__ElvUI__, __ElvUI_Libraries__, and __ElvUI_Options__) in the AddOns directory.

## Dependencies

This program uses the following Go standard library packages:

  - archive/zip
  - encoding/json
  - fmt
  - io
  - net/http
  - os
  - path/filepath
  - time

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/qqMelon/upgraded-spork/blob/trunk/LICENSE) file for details.

## Acknowledgments

  - [tukui-org](https://github.com/tukui-org) for providing the ElvUI addon
