# Demo Prefixer
A simple program designed to generate config files to automatically prefix demos with the name of the map they are recorded on. This is conceptually inspired by [this script on TFTV](https://www.teamfortress.tv/47180/demo-support-ds-prefix-on-any-map).  
The program's main purpose is to be a simple-ish program for me to test and compare different programming languages, currently Odin and Go. When you download a release it will be compiled from the Odin code in the repository.

### Important Notes: ###
- You must use the Demo Support command (`ds_record`) to record demos for the prefix to apply!
- Prefixes will not be applied to any maps you've downloaded after running the Demo Prefixer! You can set it to run automatically in various ways, see below.
- The program will not apply prefixes to demos retroactively! It will only prefix demos recorded after you've run the program.

# Installation
- Download the latest release from the [Releases page]().
- Put demo-prefixer.exe in your **tf/custom** folder.
- Run demo-prefixer.exe and answer any prompts it gives you.
- You'll see a new folder appear called **demo-prefixer** that will contain the generated config files.

# Automated Running
For convenience, the Demo Prefixer is designed to be able to be run with commandline arguments to allow for scripted automation.

There are 4 accepted arguments to allow the program to be run without user interaction: 
| Argument | Description |
| :-- | :-- |
| CompOnly | Only generates for Attack/Defend, Control Points, King of the Hill, or Payload. |
| CompPlus | Generates like CompOnly with the addition of BBall, PASS Time, UltiDuo, and UltiTrio.<br>Will take priority over CompOnly if both are present. |
| NoSuffix | Generates without map version suffixes (example_**a3**, example_**rc5**, etc).<br>Will keep event and other non-version number suffixes. |
| Silent | Runs the program without prompting for user input for automation purposes.<br>All prompted options will assume a "No" response unless specified with an argument. |

The easiest way to make use of these arguments for automatically running the program is with a batch script.  
To set that up:
- Create a text file (.txt) with any title. For example: `prefixed TF2.txt`.
- Open your text file with **Notepad**.
- Inside, copy and paste the following:
```
cd "[REPLACE WITH YOUR FILEPATH]\tf\custom"
start demo-prefixer.exe silent [OTHER ARGUMENT] [ANOTHER ARGUMENT]
start steam://rungameid/440
exit
```
- Rename your text file to change the extension from `.txt` to `.bat`. For example: `prefixed TF2.bat`
- Now running the batch script will run the Demo Prefixer and launch TF2 immediately after.
### Explaination of the batch script:
1. Navigates from where ever your batch file is to your custom folder.
2. Starts the Demo Prefixer with the silent argument and any additional arguments you set.
3. Starts TF2 through Steam, the same as if you launched via your Steam Library.
4. Exits the terminal window that the batch script will open when it runs.

There are other ways to schedule scripts to run on your computer. For example on computer start, daily, weekly, etc. I'm not going to list those here for brevity and the previous method will cover any non-power-user cases.

# How Does It Work?
The program collects the maps from TF2's various map storage locations then filters maps according to user input / arguments. Once filters have been applied, the program generates config files for each remaining map with the name of the map formatted according to user input / arguments.

On my computer this takes less than **30 milliseconds** to generate for roughly **300 maps**. I also tried make meaningful use of Odin's manual memory management to keep memory usage minimal, though likely unnecessary given the miniscule nature of the program.

# Is This A Virus?
### Nah. 
Your antivirus will likely warn you against running unknown executable files, which is reasonable advice. I'd encourage you to review the code in the repository for yourself and compile it from Source. While you're there, you can also tell me how bad my code is or how much I underutilized the languages!

Both the Go and Odin code are meant to be as close to functionally identical as I could get the two languages, pick whichever one you prefer. Odin will run slightly faster and hopefully more efficiently, Go is more well known and if you're reviewing my code you probably already have the Go compiler.
