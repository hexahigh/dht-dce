# dht-dce
![Build and release](https://github.com/hexahigh/dht-dce/actions/workflows/build&release.yml/badge.svg)
[![Release](https://img.shields.io/github/release/hexahigh/dht-dce.svg)](https://github.com/hexahigh/dht-dce/releases)
[![License](https://img.shields.io/github/license/hexahigh/dht-dce)](https://github.com/hexahigh/dht-dce/blob/main/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/hexahigh/dht-dce/total.svg)](https://github.com/hexahigh/dht-dce/releases)<br>
A simple program that converts Discord History Tracker exports to Discord chat exporter's json format.

## Download
You can download the latest nightly version [here.](https://github.com/hexahigh/dht-to-dce/releases/tag/latest_auto) (Recommended)<br>
Or, you can download the latest release [here.](https://github.com/hexahigh/dht-to-dce/releases/latest)

## Usage
<pre>
  -channel string
        Comma separated list of one or more channel id
  -id-as-name
        If set to true then the output json file will be the channel id. Overrides out.
  -in string
        path to the input SQLite database (default "input.db")
  -out string
        path to the output JSON file (default "output.json")
  -v    prints the version
</pre>

## Current limitations
<pre>
  Can only export one channel at a time.
  Does not support profile pictures or files.
</pre>
