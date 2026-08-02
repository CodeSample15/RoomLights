echo "Cloning rpi_ws281x source repo"

git clone https://github.com/jgarff/rpi_ws281x.git

cd rpi_ws281x
mkdir build && cd build
cmake -D WIDTH=300 -D HEIGHT=1 ..
cmake --build .
sudo make install
