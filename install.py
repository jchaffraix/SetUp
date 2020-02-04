import os
import platform
import subprocess
import sys

# TODO: Add flags for these.
# Configurable variables.
# TODO: Windows?
HOME=os.path.expanduser('~')
SETUP_PATH='SetUp'
DEBUG = True

# Non-configurable variables.
REMOTE='git@github.com:jchaffraix/SetUp.git'
DEPS=['git', 'tmux', 'zsh']

def _RunCommand(command):
    if DEBUG:
      print("Running command %s" % command)
    process = subprocess.Popen(command, shell=False, stdout=subprocess.PIPE)
    process.wait()
    return (process.returncode, process.stdout.read().strip())

def install_software_deps():
  print("✨ Installing deps")
  platform_os = platform.system()
  if platform_os == 'Linux':
    # We only support Debian based packet managers.
    _RunCommand(['sudo', 'apt-get', 'install'] + DEPS)
  elif platform_os == 'Darwin':
    # Brew install the previous ones.
    pass
  elif platform_os == 'Windows':
    print("Can't install missing deps on Windows.")
  else:
    print("Unknown OS %s" % platform_os)

def install_config_file(path, tools, config):
  # Check if the path exists.
  name = config.split("/")[-1]
  # TODO: This is inspired by Linux and won't work
  # on Windows.
  dst = path + "/." + name
  print("📝 Installing file %s" % dst)
  if os.path.exists(dst):
    while True:
      print("File exist %s: Overwrite/Skip/Exit [ose]: " % dst)
      answer = input()
      if answer in ["e", "E"]:
        print("💥 Exiting...")
        sys.exit(-1)
      if answer in ["s", "S"]:
        print("🚨 Skipping %s" % dst)
        return
      if answer in ["o", "O"]:
        os.rename(dst, dst + ".bak")
        break

      print("🚩 Unknown input. Try again")

  os.symlink(os.path.join(tools, config), dst)

def install_config():
  print("⚙️  Cloning the configs")
  # Create the path and git clone into it.
  clone_path = os.path.join(HOME, SETUP_PATH)
  #os.makedirs(PATH, exist_ok=True)
  # TODO: Check if this worked!
  _RunCommand(['git', 'clone', REMOTE, clone_path])

  print("🚀 Installing the configs")
  # TODO: Windows.
  install_config_file(HOME, SETUP_PATH, "/".join(["Configs", "git", "gitconfig"]))

def install():
  install_software_deps()
  install_config()
  # TODO: We don't look for failures anywhere...
  print("✅ Install successful")

if __name__ == '__main__':
  install()
