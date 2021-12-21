# cbonsai

<img src="https://i.imgur.com/rnqJx3P.gif" align="right" width="400px">

`cbonsai` is a bonsai tree generator, written in `C` using `ncurses`. It intelligently creates, colors, and positions a bonsai tree, and is entirely configurable via CLI options-- see [usage](#usage). There are 2 modes of operation: `static` (see finished bonsai tree), and `live` (see growth step-by-step).

`cbonsai` is always looking for ideas for improvement- feel free to open an issue if you've got an idea or a bug!

<br>
<br>
<br>
<br>
<br>
<br>

## Installation

<a href="https://repology.org/project/cbonsai/versions">
    <img src="https://repology.org/badge/vertical-allrepos/cbonsai.svg" alt="Packaging status" align="right">
</a>

`cbonsai` is available in multiple repositories. Check the repology chart to the right to see if `cbonsai` is packaged for your system. A big thank you to all the people who packaged `cbonsai`!

If no package exists for your system/distribution, you'll have to use the [manual](https://gitlab.com/jallbrit/cbonsai#manual) install instructions. Below are some more specific instructions for some distributions.

### Debian-based

`cbonsai` is available in Debian Testing and Unstable via `apt`. Robin Gustafsson has also kindly packaged `cbonsai` as a `.deb` file over in [this repository](https://gitlab.com/rgson/debian_cbonsai/-/packages).

### Fedora

Mohammad Kefah has kindly packaged `cbonsai` in the [Fedora copr](https://copr.fedorainfracloud.org/), which is "similar to what the AUR is to Arch". On Fedora, it can be installed like so:

```bash
sudo dnf copr enable keefle/cbonsai
sudo dnf install cbonsai
```

### MacOS

Follow the [Manual](#manual) installation, but if you install `ncurses` via homebrew, you may see this:

```
For pkg-config to find ncurses you may need to set:
  set -gx PKG_CONFIG_PATH "/usr/local/opt/ncurses/lib/pkgconfig"
```

You may need to follow these instructions before running `make install`.

If you are having trouble installing on MacOS, try reading [this issue](https://gitlab.com/jallbrit/cbonsai/-/issues/10).

#### MacPorts

On macOS, you may also install `cbonsai` using [MacPorts](https://www.macports.org). Simply install MacPorts, then issue the following commands:

```bash
sudo port selfupdate
sudo port install cbonsai
```

### Manual

You'll need to have a working `ncursesw` library. If you're on a `Debian`-based system, you can install `ncursesw` like so:

```bash
sudo apt install libncursesw5-dev
```

Or on Fedora:

```bash
sudo dnf install ncursesw5-devel
```

Once dependencies are met, then install:

```bash
git clone https://gitlab.com/jallbrit/cbonsai
cd cbonsai

# install for this user
make install PREFIX=~/.local

# install for all users
sudo make install
```

## Usage

```
Usage: cbonsai [OPTION]...

cbonsai is a beautifully random bonsai tree generator.

Options:
  -l, --live             Live mode: show each step of growth
  -t, --time=TIME        In live mode, wait TIME secs between
                           steps of growth (must be larger than 0) [default: 0.03]
  -i, --infinite         Infinite mode: keep growing trees
  -w, --wait=TIME        In infinite mode, wait TIME between each tree
                           generation [default: 4.00]
  -S, --screensaver      Screensaver mode; equivalent to -li and
                           quit on any keypress
  -m, --message=STR      Attach message next to the tree
  -T, --textOrigin=y,x   Display text from STDIN at row Y, column X
  -b, --base=INT         Ascii-art plant base to use, 0 is none
  -y, --baseY=INT        Row of the upper-left corner of the plant base
  -x, --baseX=INT        Column of the upper-left corner of the plant base
  -c, --leaf=LIST        List of comma-delimited strings randomly chosen
                           for leaves
  -M, --multiplier=INT   Branch multiplier; higher -> more
                           branching (0-20) [default: 5]
  -L, --life=INT         Life; higher -> more growth (0-200) [default: 32]
  -p, --print            Print tree to terminal when finished
  -s, --seed=INT         Seed random number generator
  -W, --save=FILE        Save progress to file [default: $XDG_CACHE_HOME/cbonsai or $HOME/.cache/cbonsai]
  -C, --load=FILE        Load progress from file [default: $XDG_CACHE_HOME/cbonsai]
  -v, --verbose          Increase output verbosity
  -h, --help             Show help
```

## Tips

### Screensaver Mode

Try out `-S/--screensaver` mode! As the help message states, it activates the `--live` and `--infinite` modes, quits upon any keypress, also saves/loads using the default cache file (`~/.cache/cbonsai`). This means:

* When you start `cbonsai` with `--screensaver`, a tree (including its seed and progress) is loaded from the default cache file.
* When you quit `cbonsai` and `--screensaver` was on, the current tree being generated (including its seed and progress) is written to the default cache file.

This is helpful for a situations like the following: let's say you're growing a really big tree, really slowly:

```bash
$ cbonsai --life 40 --multiplier 5 --time 20 --screensaver
```

Normally, when you quite `cbonsai` (e.g. by you hitting `q` or `ctrl-c`), you would lose all progress on that tree. However, by specifying `--screensaver`, the tree is automatically saved to a cache file upon quitting. The next time you run that exact same screensaver command:

```bash
$ cbonsai --life 40 --multiplier 5 --time 20 --screensaver
```

The tree is automatically loaded from the cache file! And, since infinite mode is automatically turned on, it will finish the cached tree and just keep generating more. When you quit `cbonsai` again, the tree is once again written to the cache file for next time.

Keep in mind that only the seed and number of branches are written to the cache file, so if you want to continue a previously generated tree, make sure you re-specify any other options you may have changed.

### Add to `.bashrc`

For a new bonsai tree every time you open a terminal, just add the following to the end of your `~/.bashrc`:

```bash
cbonsai -p
```

Notice it uses the print mode, so that you can immediately start typing commands below the bonsai tree.

It's also possible place the bonsai next to some text from STDIN. For example:
```bash
# Get and format the text
t=$(./cbonsai -h | sed 's/^/    /g')

# Show the text with a bonsai next to it at the bottom of the terminal window.
echo "$t" | ./cbonsai -p -T $((LINES - $(echo -n "$t" | wc -l))),0 -x 120
```

## How it Works

`cbonsai` starts by drawing the base onto the screen, which is basically just a static string of characters. To generate the actual tree, `cbonsai` uses a ~~bunch of if statements~~ homemade algorithm to decide how the tree should grow every step. Shoots to the left and right are generated as the main trunk grows. As any branch dies, it branches out into a bunch of leaves.

`cbonsai` has rules for which character and color it should use for each tiny branch piece, depending on things like what type of branch it is and what direction it's facing.

The algorithm is tweaked to look best at the default size, so larger sized trees may not be as bonsai-like.

## Inspiration

This project wouldn't be here if it weren't for its *roots*! `cbonsai` is a newer version of [bonsai.sh](https://gitlab.com/jallbrit/bonsai.sh), which was written in `bash` and was itself a port of [this bonsai tree generator](https://avelican.github.io/bonsai/) written in `javascript`.
