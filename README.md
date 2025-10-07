# Demo Prefixer
A simple program designed to automatically replace the prefix of a recorded demo file with the map name the demo was recorded on.  
The program's main purpose is to be a simple-ish program for me to test and compare different programming languages, currently just Go.

# Installation & Usage
- Download the latest release from the [Releases page](https://github.com/PapaPeach/demo-prefixer/releases).
- Put prefixer.exe in a folder containing demo (.dem) files you'd like to prefix.
- Run prefixer.exe and answer any prompts it gives you.

On my computer it runs at roughly 100 demos renamed per second.

# How Does It Work?
The [demo file format](https://github.com/mastercomfig/tf2-patches-old/blob/744d6eb003c7598f39bfdbf9046109f9d85cf11b/src/public/demofile/demoformat.h#L48) specifies a file header which contains of a couple pieces of data, including the name of the map the demo was recorded on. This header is *presumably* a rigid format, so by jumping to the specific section containing the map name, we can extract it as plain text.

# Is This A Virus?
### Nah.
Your antivirus will likely warn you against running unknown executable files, which is its job. I'd encourage you to review the code in the repository for yourself and compile it from Source. While you're there, you can also tell me how bad my code is or how much I underutilized the languages :)
