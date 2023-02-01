# RPS Arena 🪨 📄 ✂️

RPS Arena is a zero player game. The game plays itself. You watch.

<https://user-images.githubusercontent.com/8356936/212448515-22d69868-4c22-4f08-957a-99bbafe2004f.mp4>

## Installation

### Binary

Download the binary from <https://github.com/tom-on-the-internet/rps-arena/releases/latest>.

Make sure you choose the correct binary for your system.

### nix

There's a flake in this repo. So, you can Run

```sh
nix shell github:tom-on-the-internet/rps-arena
```

### Compile

Clone the repository, `go get`, `go build .`.

### Go Install

```sh
go install github.com/tom-on-the-internet/rps-arena@main
rps-arena
```

### Docker

```sh
docker run --rm -it golang bash -c 'go install github.com/tom-on-the-internet/rps-arena@main && rps-arena'
```

## Usage

Execute the binary. Press `h` for help. Enjoy!

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

## Contact

Twitter: [@tom_on_the_net](https://twitter.com/tom_on_the_net)

Email: tom@tomontheinternet.com

Project Link: [https://github.com/tom-on-the-internet/rps-arena](https://github.com/tom-on-the-internet/rps-arena)

## Acknowledgments

- Steven for inspiring me.
- [Charm](https://github.com/charmbracelet)
