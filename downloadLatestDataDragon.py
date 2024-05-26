import urllib.request, tarfile
import requests
import json
import os, shutil

datadir = "./lolData/"

def dalldir():
    for filename in os.listdir(datadir):
        file_path = os.path.join(datadir, filename)
        try:
            if os.path.isfile(file_path) or os.path.islink(file_path):
                os.unlink(file_path)
            elif os.path.isdir(file_path):
                shutil.rmtree(file_path)
        except Exception as e:
            print(f'Failed to delete {file_path}. Reason: {e}')

def extract(tar_url, extract_path='.'):
    print(tar_url)
    tar = tarfile.open(tar_url, 'r')
    for item in tar:
        tar.extract(item, extract_path)
        if item.name.find(".tgz") != -1 or item.name.find(".tar") != -1:
            extract(item.name, "./" + item.name[:item.name.rfind('/')])

def main():
    dalldir()
    response = requests.get("https://ddragon.leagueoflegends.com/api/versions.json")
    versions = json.loads(response.content)
    v = versions[0]
    print(v)
    filename = f"dragontail-{v}.tgz"
    urllib.request.urlretrieve(f"https://ddragon.leagueoflegends.com/cdn/{filename}", datadir + filename)
    extract(datadir + filename, datadir + "untarLoLData")
    with open(datadir + "untarLoLData/version.json", "w") as f:
        f.write(f"\"{v}\"")

if __name__ == "__main__":
    main()
