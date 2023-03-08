# gobonsai

A bonsai tree generator ([cbonsai](https://gitlab.com/jallbrit/cbonsai) ported to go)

<p align="center">
    <a href="https://asciinema.org/a/V6Di0zHxIKxU8T43fFSrMUSBj">
        <img alt="screenshot" src="https://user-images.githubusercontent.com/48992448/223872544-47b4255c-79d3-4352-9cac-4ac677aed013.png">
    </a>
</p>

---

```
Usage:
  gobonsai [flags]

Examples:
  gobonsai -p --seed 42
  gobonsai -l -L 48 -M 3
  gobonsai --msg "hi" --msg-y 20
  gobonsai -S -c "&,@,§,$,%,☘️,🌿,🍎,💚,🟢,🟩"

Flags:
  -l, --live             live mode: show each step of growth
  -t, --time duration    in live mode, delay between steps of growth (default 33ms)
  -i, --infinite         infinite mode: keep growing trees
  -w, --wait duration    in infinite mode, delay between each tree (default 4s)
  -S, --screensaver      screensaver mode: equivalent to -li and quit on any keypress
  -b, --base int         base pot: big=1 small=2 (default 1)
      --base-x int       column position of upper-left corner of plant base pot
      --base-y int       row position of upper-left corner of plant base pot
  -a, --align int        align tree: center=1 left=2 right=3 (default 1)
  -m, --msg string       attach message next to the tree
      --msg-x int        column position of upper-left corner of message text (default 4)
      --msg-y int        row position of upper-left corner of message text (default 2)
  -c, --leaves string    list of comma-delimited leaves (default "&")
  -M, --multiplier int   branch multiplier higher -> more branching (0-20) (default 5)
  -L, --life int         life higher -> more growth (0-127) (default 32)
  -p, --print            print first tree to stdout and exit immediately
  -s, --seed int         seed random number generator (default random)
  -h, --help             show help
```

---

[gobonsai](https://github.com/nothub/gobonsai) is a `go` port of  
[cbonsai](https://gitlab.com/jallbrit/cbonsai) (written in `C` and itself a port of  
[bonsai.sh](https://gitlab.com/jallbrit/bonsai.sh) (written in `bash` and itself a port of  
[bonsai](https://avelican.github.io/bonsai/) (written in `javascript`))).
