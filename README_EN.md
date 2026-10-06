# AuraClock

[简体中文](README.md) | [English](README_EN.md)

## Introduction

Humans do not naturally operate according to standardized clock time. Instead, our biological rhythms are dynamically influenced by environmental changes such as sunrise and sunset. AuraClock automatically obtains sunrise and sunset times based on latitude and longitude, then calculates the time points of daily activities according to the time offsets configured by the user. This is intended to make daily schedules better aligned with the biological clock and improve working performance.

## Running

This software uses the cross-platform GUI framework Fyne, so it can be compiled and run locally. An executable for Linux X86_64 is provided in the `target` directory. If your environment is compatible, you can download and run it directly.

Try running:

```shell
go run .
```

Build a binary:

```shell
go build .
```

When running or compiling for the first time, Go will prepare the build environment, mainly by downloading and compiling third-party libraries. Please wait patiently.

## Usage

The software runs in GUI mode by default. You can also use command-line options for quick operations. All alarm data is saved to `GUI.json` in the current directory.

### GUI

The interface is straightforward:

![Example.png](readme/示例.png)

* On the left, you can set the name, time offset (in minutes), and notes for a new alarm. Click **Add Alarm** to create it.
* On the right, all alarms currently stored in `GUI.json` are displayed. Click the gray button to delete an alarm.

All changes are written to `GUI.json` when the software exits.

### CLI Options

* Option `a`

  * Input: floating-point number
  * Purpose: Set the latitude. Beijing is used by default. `(default 39.906)`

* Option `l`

  * Input: floating-point number
  * Purpose: Set the longitude. Beijing is used by default. `(default 116.391)`

* Option `n`

  * Purpose: Use only the basic alarm creation functionality.

* Option `r`

  * Input: file path (string)
  * Purpose: Specify the file from which alarms should be read.

* Option `w`

  * Input: file path (string)
  * Purpose: Create alarms and write them to the specified file.

## TODO

* Improve the GUI:

  * The delete button does not have a text label.
  * Entering the time offset directly as a number is somewhat inconvenient; a selector could be added.
* Add GUI functionality:

  * Allow users to select the file to read from and write to.
* User configuration file:

  * Store preferences and settings separately, such as latitude and longitude.

