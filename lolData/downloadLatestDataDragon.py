import urllib.request, tarfile
import requests
import json

def extract(tar_url, extract_path='.'):
    print(tar_url)
    tar = tarfile.open(tar_url, 'r')
    for item in tar:
        tar.extract(item, extract_path)
        if item.name.find(".tgz") != -1 or item.name.find(".tar") != -1:
            extract(item.name, "./" + item.name[:item.name.rfind('/')])

def main():
    response = requests.get("https://ddragon.leagueoflegends.com/api/versions.json")
    versions = json.loads(response.content)
    v = versions[0]
    print(v)
    filename = f"dragontail-{v}.tgz"
    urllib.request.urlretrieve(f"https://ddragon.leagueoflegends.com/cdn/{filename}", filename)
    extract(filename, "./untarLoLData")
    with open("./untarLoLData/version.json", "w") as f:
        f.write(f"\"{v}\"")

if __name__ == "__main__":
    main()
