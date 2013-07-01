Blinky
======

This is the source code for driving a BlinkM
based on Mixpanel queries. It's currently
driving a lamp in the Mixpanel offices.

It runs on Linux, and depends on the
[I2CTools](http://www.lm-sensors.org/wiki/I2CTools)
package. If you happen to be using a Raspberry Pi,
there is more information about setting up your
Pi to use the package here:

    http://learn.adafruit.com/adafruits-raspberry-pi-lesson-4-gpio-setup/configuring-i2c

Once your i2c is set up and your BlinkM is all
plugged in and happy, you can use blinky by
by copying the config.json.example file in the root
blinky directory, editing it with your Mixpanel project
information and Mixpanel queries, and then starting
it on the command line- something like:

    $ cd blinky
    $ go build # build the blinky executable
    $ cp config.json.example config.json
    $ vi config.json # Add your Mixpanel info
    $ ./blinky config.json
