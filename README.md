# Opus Packet Decoder
This project utilizes a golang wrapper for the libopus library in order decode audio data and obtain their pcm data.

## Instalation
The key part to this project is the libopus C library. The wrapper for this library can be found here: https://github.com/hraban/opus. More specificaly, there are two libraries that are needed; `libopus-dev` and `libopusfile-dev`. To install these you can run the following commands:

Linux:
    sudo apt-get install pkg-config libopus-dev libopusfile-dev

Mac:
    brew install pkg-config opus opusfile

## Expected data format
Right now there is only one format excepted. You pass in a file with the `-f` flag that contains the opus data packets, encoded as base64, with each packet being seperated by a new line. For example, something like this:


    SINOovDvVb3x9wGCf1RrPXd9I0ssAfgCiUA= 
    SDJMgV+WfQSe4d4= 
    SAwZydbIYEtFwvaM6A== 
    SAo77RiblygKyNtOnTXk 
    SAu6NdbJwjlwtcZwH3FO 
    SAwARKRgS1JR/4R9mZ86JWA= 
    SAlORGiFIyuiqhhIf0xVFNYs 
    SAbjecVl6VYXXLA= 

You can also specify the number of channels with the `-channel` flag and sample rate with the `-sr` flag.

## Working with ffmpeg
Once you have the pcm raw data of the recording, you can run the following command to generate a wav file:

    ffmpeg -f s16le -ar 16k -i out.pcm file.wav

The s16le represends the way the pcm raw data is encoded: signed 16 bit little endian. Chose the correct sample rate as well.
