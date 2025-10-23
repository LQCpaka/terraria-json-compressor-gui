
<p align=center>
  <img src="https://github.com/LQCpaka/terraria-json-compressor-gui/blob/main/frontend/public/images/terria-icon-logo.png"/>
<p>
    
<p align="center">


  <img src="https://img.shields.io/badge/JSON-Converter-blue">
  <img src="https://img.shields.io/badge/Vietnam-⭐_Vietnamese-red">
  <img src="https://img.shields.io/badge/Languuage-Golang-blue">
  <img src="https://img.shields.io/badge/Application-GUI-blue">
  <img src="https://img.shields.io/badge/Terraria-JSON-red">

</p>

# TERRARIA JSON COMPRESSOR GUI

```TJC``` is made by LQC - LE QUOC CAN.

Application made by LQC. For translating purpose. Help people could convert the file from CSV file (comma-separated values) to JSON file, which is working for modding stuff.

App made with Wails and Golang. With webview architect, made the app could build and run on cross-platform. And React for frontend.

## OPERATING SYSTEM - SUPPORT
<table>
  <tr>
    <th>Operating System</th>
    <th>Supported</th>
    <th>Note</th>
  </tr>
  <tr>
    <td align=center>Windows</td>
    <td align=center>✓</td>
    <td >Not sure if the application could work fine on all version of windows. But I think its fine.</td>
  </tr>
  <tr>
    <td align=center>Linux</td>
    <td align=center>✓</td>
    <td >Currently dont have installer for linux, you have to build it yourself, clone project and build it on linux.</td>
  </tr>
  <tr>
    <td align=center>Mac OS</td>
    <td align=center>✓</td>
    <td >Currently dont have installer for Mac, you have to build it yourself, clone project and build it on Mac</td>
  </tr>
</table>

## Installations - Builds 🛠️

> [!NOTE]
> **Some ways could be outdated. You could fix it yourself or post on issue tab in this repo. I currently only release for windows only, so in future, If I have time, imma gonna do a release on Linux and Mac aswell**.

___

### Linux 🐧

**1. Install Golang programming language. Nodejs - For Frontend support.**

Project required you should have at least ≥ 1.20. Should be fine if you not choosing specific version, because the install gonna take the latest version for you (above .20). Use your package manager like ```apt``` ```dnf``` ```pacman```, etc... to install support package.

```
sudo apt install golang-go
```

**2. Install WebKitGTK**

```
sudo apt install libwebkit2gtk-4.0-dev
```

Optional: You could install these things as well incase your OS doesnt suitable for it. You could review it and remove something you dont really need.

```
sudo apt install build-essential libgtk-3-dev libglib2.0-dev libgdk-pixbuf2.0-dev libpango1.0-dev libcairo2-dev
```

**3. Install Wails CLI.**

After done install golang config. You should do this. Install wails framework.

```
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**4. Build app and use.**

```cd``` to project folder and type this. After done, you could find your app after build.

```
wails build
```
___
### MacOS 🍎

**1. Install Golang programming language. Nodejs - For Frontend support.**

Project required you should have at least ≥ 1.20. Should be fine if you not choosing specific version, because the install gonna take the latest version for you (above .20). You should install on offical website golang or you could install it with ```brew```.

```
brew install go
brew install nodejs
```

**2. Install Wails-CLI**

After done install golang config. You should do this. Install wails framework.

```
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**5. Build app and use.**

```cd``` to project folder and type this. After done, you could find your app after build.

```
wails build
```
Optional: Universal build for both Intel & Apple Silicon
```
wails build -platform darwin/universal
```
## How To Use ❓

> [!NOTE]
> You could find csv and the standard of csv editing for terraria in [Terraria Forum](https://forums.terraria.org/index.php?threads/the-ultimate-guide-to-content-creation-and-use-for-the-terraria-workshop.100652/#languagepack). Also thank you everyone that contributed for that thread.

- Browse your csv file and select it.
- You could preview your file if you want make sure there is nothing go wrong with error syntax in your file.
- Press ```Start Compress``` and wait the application build for you. After that you will receive a json file. Paste it into localiztion folder of your mod folder and try it.
<br>

> [!CAUTION]
> **DO NOT** Remove header line of csv. Its could harm the process, unless you know what you are doing. What is header line? - Its is ```Key,Translation```. Example: ```Key, En-US```. A file should have only 2 columns, 1 for ```Key```, 1 for ```Translation```. Keep this rule: **If its fine, don't touch it**.

<br>

> [!TIP]
> You could use application like excel, sheet or something like that. That helps you editing the file easier. Other methods like using ```notepad++``` or something like that also work fine
