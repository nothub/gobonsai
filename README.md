# gobonsai

A bonsai tree generator, written in using [tcell](https://github.com/gdamore/tcell).

## Flags

```
  -l, --live             live mode: show each step of growth
  -t, --time duration    in live mode, wait TIME between steps of growth (default 33ms)
  -i, --infinite         infinite mode: keep growing trees
  -w, --wait duration    in infinite mode, wait TIME between each tree generation (default 4s)
  -S, --screensaver      screensaver mode: equivalent to -li and quit on any keypress
  -b, --base int         base pot: big=1 small=2 (default 1)
  -c, --leaves string    list of comma-delimited strings randomly chosen for leaves (default "&")
  -M, --multiplier int   branch multiplier higher -> more branching (0-20) (default 5)
  -L, --life int         life higher -> more growth (0-200) (default 32)
  -p, --print            print tree to stdout on exit
  -s, --seed int         seed random number generator (default 42)
  -h, --help             show help
  -m, --message string   attach message next to the tree
  -a, --align int        align tree: center=0 left=1 right=2

```

---

This project wouldn't be here if it weren't for its *roots*! [gobonsai](](https://gitlab.com/nothub/gobonsai)) is a port
of [cbonsai](https://gitlab.com/jallbrit/cbonsai), which was written in `C` and was itself a port
of [bonsai.sh](https://gitlab.com/jallbrit/bonsai.sh), which was written in `bash` and was itself a port
of [bonsai](https://avelican.github.io/bonsai/), which was written in `javascript`.
