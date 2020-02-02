import platform
import subprocess
import os

# TODO: Add flags for these.
# Configurable variables.
# TODO: Windows?
PATH=os.path.expanduser('~/Tools')
DEBUG = True

# Non-configurable variables.
REMOTE='git@github.com:jchaffraix/Tools.git'
DEPS=['git', 'tmux', 'zsh']

def _RunCommand(command):
    if DEBUG:
      print("Running command %s" % command)
    process = subprocess.Popen(command, shell=False, stdout=subprocess.PIPE)
    process.wait()
    return (process.returncode, process.stdout.read().strip())

def install_software_deps():
  print("Installing deps")
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

def install_config():
  print("Installing configs")
  # Create the path and git clone into it.
  os.makedirs(PATH, exist_ok=True)
  _RunCommand(['git', 'clone', REMOTE, PATH])


def install():
  install_software_deps()
  install_config()

if __name__ == '__main__':
  install()
