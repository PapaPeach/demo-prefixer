# Demo Prefixer
A simple program designed to generate config files to automatically prefix demos with the name of the map they are recorded on. This is conceptually inspired by [this script on TFTV](https://www.teamfortress.tv/47180/demo-support-ds-prefix-on-any-map).

### Important Notes: ###
- You must use the Demo Support command (`ds_record`) to record demos for the prefix to apply!
- Prefixes will not be applied to any maps you've downloaded after running the Demo Prefixer! You can set it to run automatically in various ways, see below.
- The program will not apply prefixes to demos retroactively! It will only prefix demos recorded after you've run the program.

# Installation
- Download the latest release from the [Releases page]().
- Put demo-prefixer.exe in your **tf/custom** folder.
- Run demo-prefixer.exe and answer any prompts it gives you.

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

# How Does It Work?

# Is This A Virus?
