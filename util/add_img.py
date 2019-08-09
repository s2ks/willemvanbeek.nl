#!/bin/python

import glob
import sqlite3
import sys
import subprocess

conn = sqlite3.connect(sys.argv[1])
c = conn.cursor()
desc = """
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. A erat nam at lectus urna duis convallis convallis.
"""

print(sys.argv)

for path in glob.glob(sys.argv[2] + '*'):
    filename = path.split('/')
    filename = filename[len(filename) - 1]
    subprocess.run(["feh", path])
    c.execute("INSERT INTO beelden (name, material, description, img_path) VALUES(?, ?, ?, ?);", (input("Name: "),
        input("Material (steen, hout, metaal): "), desc, "/static/img/" + filename))

conn.commit()
conn.close()
