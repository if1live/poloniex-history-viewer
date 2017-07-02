#!/usr/bin/env python3

import requests
import re

def read_asset_list(filename):
    f = open(filename)

    #<link rel="stylesheet" type="text/css" href="/css/jquery.qtip.min.css?v=122015" />
    #<script src="/js/jquery-2.1.3.min.js" type="text/javascript"></script>
    prog_css = re.compile(r"href=(?P<path>.+\.css)")
    prog_js = re.compile(r"src=(?P<path>.+\.js)")

    css_list = []
    js_list = []

    for line in f:
        # css
        m = prog_css.search(line)
        if m:
            p = m.group('path')
            p = p.replace("'", '').replace('"', "")
            css_list.append(p)


        # js
        m = prog_js.search(line)
        if m:
            p = m.group('path')
            p = p.replace("'", '').replace('"', "")
            js_list.append(p)

    return css_list + js_list

def download_asset(path):
    url = 'https://poloniex.com' + path
    r = requests.get(url)
    text = r.text

    f = open("." + path, "w")
    f.write(text)
    f.close()
    print("download: " + url)


assets = read_asset_list('trade_history.html')
for a in assets:
    download_asset(a)


download_asset('/css/fonts/Roboto/Roboto-Light.ttf')
download_asset('/css/fonts/Roboto/Roboto-Regular.ttf')
download_asset('/css/fonts/Roboto/Roboto-Bold.ttf')

download_asset('/css/fonts/fontawesome/fontawesome-webfont.eot?v=4.2.0')
download_asset('/css/fonts/fontawesome/fontawesome-webfont.woff?v=4.2.0')
download_asset('/css/fonts/fontawesome/fontawesome-webfont.ttf?v=4.2.0')
download_asset('/css/fonts/fontawesome/fontawesome-webfont.svg?v=4.2.0#fontawesomeregular')
