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

for filename in glob.glob(sys.argv[2] + '*'):
    subprocess.run(["feh", filename])
    c.execute("INSERT INTO beelden (name, material, description, img_path) VALUES(?, ?, ?, ?);", (input("Name: "),
        input("Material (steen, hout, metaal): "), desc, filename))

conn.commit()
conn.close()
