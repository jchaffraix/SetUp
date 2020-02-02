import platform
import subprocess

# TODO: Add flags for these.
# Configurable variables.
# TODO: Windows?
PATH="~/Tools"
DEBUG = True

# Non-configurable variables.
#REMOTE="https://github.com/jchaffraix/Tools.git
DEPS=['git', 'tmux', 'zsh']

def _RunCommand(command):
    if DEBUG:
      print "Running command %s" % command
    process = subprocess.Popen(command, shell=False, stdout=subprocess.PIPE)
    process.wait()
    return (process.returncode, process.stdout.read().strip())

def install_software_deps():
  os = platform.system()
  if os == 'Linux':
    # We only support Debian based packet managers.
    print("Installing deps")
    _RunCommand(['sudo', 'apt-get', 'install'] + DEPS)
  elif os == 'Darwin':
    # Brew install the previous ones.
    pass
  elif os == 'Windows':
    print("Can't install missing deps on Windows.")
  else:
    print("Unknown OS %s" % os)

def install():
  install_software_deps()

if __name__ == '__main__':
  install()
