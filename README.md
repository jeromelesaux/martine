# martine

After Claudia from Eliot a long time ago, here is coming Martine.
A cli which converts JPEG or PNG image file into  SCR/PAL or Overscan  Amstrad CPC file screen.
Multi os, you can convert any pictures to Amstrad CPC Screen.
The files generated (.win, .scr, .ink) are compatible with [OCP art studio](http://www.cpc-power.com/index.php?page=detail&num=4963) and Impdraw V2 [i2](http://amstradplus.forumforever.com/t462-iMPdraw-v2-0.htm)

# Contents :
 * [Introduction](#Introduction)
 * [Output files](#Output_files)
 * [Options](#Options)
 * [Installation](#Installation)
 * [Usage](#Usage)
 * [Technics](#Technics)
	* [Screen Conversion](#Screen_conversion)
	* [Roll sprite](#Roll)
	* [Tiles](#Tiles)
	* [Rotation Sprite](#Rotation)
	* [3D Rotation](#3D_rotation)
	* [Tilemap](#Tilemap)
	* [Flash](#flash)
	* [Egx](#Egx)
	* [Deltapacking](#Deltapacking)

## Introduction 
Martine tries to accelerate your game, demo animation development by organize and conversion of your graphical data.
All screens data are associated with basic programs launchers to quick display the output on real machines.
You can follow the development here [Cpcwiki in english](https://www.cpcwiki.eu/forum/applications/martine-a-cpc-old-plus-tool-to-handle-your-images/) or here [Amstrad plus forum in french](https://amstradplus.forumforever.com/t513-martine-fait-du-dessin.htm) <br>
A Gui was developed for windows users by Tronic you can get it [here](https://index.amstrad.info/wp-content/uploads/martine/)



## Output_files
Output:
Martine generates a lot of differents files, which allow you to find the format and the data you will use.

### files generated :  
 * .win or .scr sprite or screen files
 * .pal or .ink palette file (.ink will be generated if the -p option is set)
 * .txt ascii file with palettes values (firmware values and basic values), and screen byte values
 *  c.txt ascii file with palettes values (firmware values and basic values), and screen byte values by column
 * .json json file with palettes values (firmware values and basic values), and screen byte values
 * .bas launch to test the screen load on classic .scr 17ko 
 * _resized.png images files to ensure the resize action
 * _downgraded.png images files to ensure the downgraded palette action
 * _paletteink.png palette image describing the CPC old colors used
 * _palettepal.png palette image describing the CPC old colors used
 * _palettekit.png palette image describing the CPC plus colors used

## Options 
### additionnals options available : 
* -dsk will generate a dsk file and add all amsdos files will be added.
* -n will remove amsdos headers from the amsdos files
* -f will generate overscan screen amsdos file
* -p will generate a CPC plus screen amsdos file
* -h will generate sprite of x pixel high
* -w will generate sprite of x pixel wide
* -s to define the byte token will be replace the byte token in the ascii files
* -a to set the algorithm to downsize the image
* -m to define the screen mode 0,1,2
* -o to set the output directory (is not exists, martine will create the folder).
* -z compress the image .scr / .win (not overscan) using the algorithm.
* -zigzag will generate data output as zigzag.
* -scanlinesequence will generate the lines following the input sequence such as 0,7,1,2,3,4,5,6 for instance.
* -sla will rotate x column pixels from the left, those columns will be discarded  
* -sra will rotate x column pixels from the right, those columns will be discarded  
* -rra will rotate x column pixels from the right
* -rla will rotate x column pixels from the left
* -keephigh will rotate x line pixels to the top
* -keeplow will rotate x line pixels to the bottom
* -losthigh will rotate x line pixels to the top, those lines will be discarded 
* -lostlow will rotate x line pixels to the bottom, those lines will be discarded
* -sna copy output files in a new CPC image Sna.
* -reducer reducing color filter (3 gradients are available)
* -mask string
    	Mask to apply on each bit of the sprite (to apply an and operation on each pixel with the value #AA [in hexdecimal: #AA or 0xAA, in decimal: 170] ex: martine -i myimage.png -w 40 -h 80 -mask #AA -m 0 -maskand)
* -maskand 	Will apply an AND operation on each byte with the mask
* -maskor Will apply an OR operation on each byte with the mask
* -tilemap
    	Analyse the input image and generate the tiles, the tile map and gloabl schema.
* -spritehard will generate 16x16 bits sprite hard for CPC plus.
* -splitrasters will generate a rastered screen 
* -reverse create a png image from your .scr .win file

### hardware options (if you owns a M4 Card, you can transfert your results by Wifi to your CPC using those options) : 
* -host ip or dns name of your M4.
* -autoexec will execute the launcher or sna file on your remote CPC.


## Installation

To Install and compile
```
go get github.com/jeromelesaux/martine
cd $GOPATH/src/github.com/jeromelesaux/martine
go get 
go build
```

To get binary : 
[https://github.com/jeromelesaux/martine/releases](https://github.com/jeromelesaux/martine/releases)
<br>OS avaible : Linux, Macos X and Windows  


## Usage
Usage and options : 

```
./martine
martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)
By Impact Sid (Version:0.25rc)
Special thanks to @Ast (for his support), @Siko and @Tronic for ideas
usage :

  -a int
    	Algorithm to resize the image (available :
    		1: NearestNeighbor (default)
    		2: CatmullRom
    		3: Lanczos
    		4: Linear
    		5: Box
    		6: Hermite
    		7: BSpline
    		8: Hamming
    		9: Hann
    		10: Gaussian
    		11: Blackman
    		12: Bartlett
    		13: Welch
    		14: Cosine
    		15: MitchellNetravali
    		 (default 1)
  -address string
    	Starting address to display sprite in delta packing (default "0xC000")
  -animate
    	Will produce an full screen with all sprite on the same image (add -i image.gif or -i *.png)
  -autoexec
    	Execute on your remote CPC the screen file or basic file.
  -delta
    	Delta mode: compute delta between two files (prefixed by the argument -df)
    		(ex: -delta -df file1.SCR -df file2.SCR -df file3.SCR).
    		(ex with wildcard: -delta -df file\?.SCR or -delta file\*.SCR
  -df value
    	scr file path to add in delta mode comparison. (wildcard accepted such as ? or * file filename.)
  -dithering int
    	Dithering algorithm to apply on input image
    	Algorithms available:
    		0: FloydSteinberg
    		1: JarvisJudiceNinke
    		2: Stucki
    		3: Atkinson
    		4: Sierra
    		5: SierraLite
    		6: Sierra3
    		7: Bayer2
    		8: Bayer3
    		9: Bayer4
    		10: Bayer8
    	 (default -1)
  -dsk
    	Copy files in a new CPC image Dsk.
  -egx1
    	Create egx 1 output cpc image overscan (option -f) or classical (mix mode 0 / 1).
    		(ex before generate two images one in mode 1 et one in mode 0
    		for instance : martine -i myimage.jpg -m 0 and martine -i myimage.jpg -m 1
    		: -egx1 -i 1.SCR -m 0 -pal 1.PAL -i2 2.SCR -o test -m2 1 -dsk)
    		or
    		(ex automatic egx from image file : -egx1 -i input.png -m 0 -o test -dsk)
  -egx2
    	Create egx 2 output cpc image overscan (option -f) or classical (mix mode 1 / 2).
    		(ex before generate two images one in mode 1 et one in mode 2
    		for instance : martine -i myimage.jpg -m 0 and martine -i myimage.jpg -m 1
    		: -egx2 -i 1.SCR -m 0 -pal 1.PAL -i2 2.SCR -o test -m2 1 -dsk)
    		or
    		(ex automatic egx from image file : -egx2 -i input.png -m 0 -o test -dsk)
  -extendeddsk
    	Export in a Extended DSK 80 tracks, 10 sectors 400 ko per face
  -f	Overscan mode (default no overscan)
  -flash
    	generate flash animation with two ocp screens.
    		(ex: -m 1 -flash -i input.png -o test -dsk)
    		or
    		(ex: -m 1 -flash -i input1.scr -pal input1.pal -m2 0 -i2 input2.scr -pal2 input2.pal -o test -dsk )
  -h int
    	Custom output height in pixels. (Will produce a sprite file .win) (default -1)
  -help
    	Display help message
  -host string
    	Set the ip of your M4.
  -i string
    	Picture path of the input file.
  -i2 string
    	Picture path of the second input file (flash mode)
  -info
    	Return the information of the file, associated with -pal and -win options
  -initprocess string
    	Create a new empty process file.
  -ink string
    	Path of the palette Cpc ink file. (Apply the input ink palette on the image)
  -iter int
    	Iterations number to walk in roll mode, or number of images to generate in rotation mode. (default -1)
  -iterx int
    	Number of tiles on a row in the input image. (default 1)
  -itery int
    	Number of tiles on a column in the input image. (default 1)
  -keephigh int
    	Bit rotation on the top and keep pixels (default -1)
  -keeplow int
    	Bit rotation on the bottom and keep pixels (default -1)
  -kit string
    	Path of the palette Cpc plus Kit file. (Apply the input kit palette on the image)
  -losthigh int
    	Bit rotation on the top and lost pixels (default -1)
  -lostlow int
    	Bit rotation on the bottom and lost pixels (default -1)
  -m int
    	Output mode to use :
    		0 for mode0
    		1 for mode1
    		2 for mode2
    		and add -f option for overscan export.
    		 (default -1)
  -m2 int
    	Output mode to use :
    		0 for mode0
    		1 for mode1
    		2 for mode2
    		mode of the second input file (flash mode) (default -1)
  -mask string
    	Mask to apply on each bit of the sprite (to apply an and operation on each pixel with the value #AA [in hexdecimal: #AA or 0xAA, in decimal: 170] ex: martine -i myimage.png -w 40 -h 80 -mask #AA -m 0 -maskand)
  -maskand
    	Will apply an AND operation on each byte with the mask
  -maskor
    	Will apply an OR operation on each byte with the mask
  -multiplier float
    	Error dithering multiplier. (default 1.18)
  -n	No amsdos header for all files (default amsdos header added).
  -o string
    	Output directory
  -p	Plus mode (means generate an image for CPC Plus Screen)
  -pal string
    	Apply the input palette to the image
  -pal2 string
    	Apply the input palette to the second image (flash mode)
  -processfile string
    	Process file path to apply.
  -quantization
    	Use additionnal quantization for dithering.
  -reducer int
    	Reducer mask will reduce original image colors. Available :
    		1 : lower
    		2 : medium
    		3 : strong
    	 (default -1)
  -remotepath string
    	Remote path on your M4 where you want to copy your files.
  -reverse
    	Transform .scr (overscan or not) file with palette (pal or kit file) into png file
  -rla int
    	Bit rotation on the left and keep pixels (default -1)
  -roll
    	Roll mode allow to walk and walk into the input file, associated with rla,rra,sra,sla, keephigh, keeplow, losthigh or lostlow options.
  -rotate
    	Allow rotation on the input image, the input image must be a square (width equals height)
  -rotate3d
    	Allow 3d rotation on the input image, the input image must be a square (width equals height)
  -rotate3dtype int
    	Rotation type :
    		1 rotate on X axis
    		2 rotate on Y axis
    		3 rotate reverse X axis
    		4 rotate left to right on Y axis
    		5 diagonal rotation on X axis
    		6 diagonal rotation on Y axis

  -rotate3dx0 int
    	X0 coordinate to apply in 3d rotation (default width of the image/2) (default -1)
  -rotate3dy0 int
    	Y0 coordinate to apply in 3d rotation (default height of the image/2) (default -1)
  -rra int
    	Bit rotation on the right and keep pixels (default -1)
  -s string
    	Byte statement to replace in ascii export (default is db), you can replace or instance by defb or byte
  -scanlinesequence string
    	Scanline sequence to apply on sprite. for instance :
    		martine -i myimage.jpg -w 4 -h 4 -scanlinesequence 0,2,1,3
    		will generate a sprite stored with lines order 0 2 1 and 3.

  -sla int
    	Bit rotation on the left and lost pixels (default -1)
  -sna
    	Copy files in a new CPC image Sna.
  -splitrasters
    	Create Split rastered image. (Will produce Overscan output file and .SPL with split rasters file)
  -spritehard
    	Generate sprite hard for cpc plus.
  -sra int
    	Bit rotation on the right and lost pixels (default -1)
  -tile
    	Tile mode to create multiples sprites from a same image.
  -tilemap
    	Analyse the input image and generate the tiles, the tile map and gloabl schema.
  -version
    	print martine's version
  -w int
    	Custom output width in pixels. (Will produce a sprite file .win) (default -1)
  -win string
    	Filepath of the ocp win file
  -z int
    	Compression algorithm :
    		1: rle (default)
    		2: rle 16bits
    		3: Lz4 Classic
    		4: Lz4 Raw
    	 (default -1)
  -zigzag
    	generate data in zigzag order (inc first line and dec next line for tiles)
```
## Technics
Principles : 
Martine can be used in several modes : 
 * first mode : conversion of an input image into sprite (file .win), cpc screen (file .scr) and overscan screen (larger file .scr)
 * second mode : rotate line or column pixels of an input image (option -roll)
 * bulk conversion of a tiles page : each sprites of the same width and height will be converted to sprites (files .win, option -tile)
 * rotate an existing image, will produce images rotated
 * 3d rotate an existing image on an axis.
 * egx mode screen generation (screen which alternates 2 modes (0,1,2) each line).
 * delta packing, allows you to create your own animation using the less memory as possible, keeping the maximum of speed.
 * tilemap generate the map of tiles and sprites of the input map.
 * split raster to generate a rastered screen
 * flash to generate a flash screen

examples :
### Screen_conversion 
* convert samples/Batman-Neal-Adams.jpg 

  * in mode 0 
```martine -i samples/Batman-Neal-Adams.jpg -m 0```
  * in mode 1 
```martine -i samples/Batman-Neal-Adams.jpg -m 1```
  * in mode 2 
```martine -i samples/Batman-Neal-Adams.jpg -m 2```
  * in mode 0 in overscan : 
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -f```
  * in mode 0 overscan for Plus series :
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -f -p```
  * to get sprites (40 pixels wide)
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -w 40```
  * roll mode to do an rra operation on the image (will create 16 sprites with a rra operation on the first pixels on the left)
	```martine -i samples/Batman-Neal-Adams.jpg -m 0 -w 40 -roll -rra 1 -iter 16```	

Samples : 
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -f```

input ![samples/Batman-Neal-Adams.jpg](samples/Batman-Neal-Adams.jpg)      
 
 will resize the image and save it 

 ![resized](samples/batman_mode0_resized.png)

 after downgrade the colors palette : 

 ![downgrade colors](samples/batman_mode0_down.png)

 results on a CPC emulator : 

 ![result](samples/overscan-batman.png)

 


### Roll : 

```martine -i samples/rotate.png -m 0 -w 16 -h 16 -roll -rra 1 -iter 16```

input ![samples/rotate.png](samples/rotate.png)

sames phasis reduce size and downgrade colors palette to CPC palette. 

after rotate the first pixels' column in 16 differents images : 

 ![0rotate.png](samples/0rotate.png)
 ![1rotate.png](samples/1rotate.png)
 ![2rotate.png](samples/2rotate.png)
 ![3rotate.png](samples/3rotate.png)
 ![4rotate.png](samples/4rotate.png)
 ![5rotate.png](samples/5rotate.png)
 ![6rotate.png](samples/6rotate.png)
 ![7rotate.png](samples/7rotate.png)
 ![8rotate.png](samples/8rotate.png)
 ![9rotate.png](samples/9rotate.png)
 ![10rotate.png](samples/10rotate.png)
 ![11rotate.png](samples/11rotate.png)
 ![12rotate.png](samples/12rotate.png)
 ![13rotate.png](samples/13rotate.png)
 ![14rotate.png](samples/14rotate.png)
 ![15rotate.png](samples/15rotate.png)

 with the same image, to rotate the pixels line : 

 ```martine -i samples/rotate.png -m 0 -w 16 -h 16 -roll -keephigh 2 -iter 16 ```

 will produce images : 

 ![0rotate.png](samples/0rotate_keephigh.png)
 ![1rotate.png](samples/1rotate_keephigh.png)
 ![2rotate.png](samples/2rotate_keephigh.png)
 ![3rotate.png](samples/3rotate_keephigh.png)
 ![4rotate.png](samples/4rotate_keephigh.png)
 ![5rotate.png](samples/5rotate_keephigh.png)
 ![6rotate.png](samples/6rotate_keephigh.png)
 ![7rotate.png](samples/7rotate_keephigh.png)
 ![8rotate.png](samples/8rotate_keephigh.png)
 ![9rotate.png](samples/9rotate_keephigh.png)
 ![10rotate.png](samples/10rotate_keephigh.png)
 ![11rotate.png](samples/11rotate_keephigh.png)
 ![12rotate.png](samples/12rotate_keephigh.png)
 ![13rotate.png](samples/13rotate_keephigh.png)
 ![14rotate.png](samples/14rotate_keephigh.png)
 ![15rotate.png](samples/15rotate_keephigh.png)

 

### Tiles :

this option will extract all the tiles from an image and generate the sprites files. 
sample usage : 
```martine -i samples/tiles.png -tile -w 64 -iterx 14 -itery 7 -m 0 ```

input ![samples/rotate.png](samples/tiles.png)

This command will generate 14*7 sprites of 64 pixels large.
Warn, all sprite must have the same size.

 ![0rotate.png](samples/TILES_resized_0.png)
 ![1rotate.png](samples/TILES_resized_1.png)
 ![2rotate.png](samples/TILES_resized_2.png)
 ![3rotate.png](samples/TILES_resized_3.png)
 ![4rotate.png](samples/TILES_resized_4.png)
 ![5rotate.png](samples/TILES_resized_5.png)
 ![6rotate.png](samples/TILES_resized_6.png)
 ![7rotate.png](samples/TILES_resized_7.png)
 ![8rotate.png](samples/TILES_resized_8.png)
 ![9rotate.png](samples/TILES_resized_9.png)
 ![10rotate.png](samples/TILES_resized_10.png)
 ![11rotate.png](samples/TILES_resized_11.png)
 ![12rotate.png](samples/TILES_resized_12.png)
 ![13rotate.png](samples/TILES_resized_13.png)
 ![14rotate.png](samples/TILES_resized_14.png)
 ![15rotate.png](samples/TILES_resized_15.png)
 ![10rotate.png](samples/TILES_resized_16.png)
 ![11rotate.png](samples/TILES_resized_17.png)
 ![12rotate.png](samples/TILES_resized_18.png)
 ![13rotate.png](samples/TILES_resized_19.png)
 ![14rotate.png](samples/TILES_resized_20.png)


### Rotation : 
This option able to rotate the input image, iter will generate the number of images.
```martine -i images/coke.jpg -rotate -iter 16 -o test  -m 1 -w 32 -h 32```

Input : ![samples/coke.jpg](samples/coke.jpg)

results : ![samples/coke.igf](samples/coke.gif)


### 3D_rotation
This option able to rotate the input image on an axis (X or Y and even on a diagonal), iter will generate the number of images. 
```./martine -m 0 -w 64 -h 64 -rotate3d -rotate3dtype 2 -o test/ -i images/sonic.png  -iter 12```

Input : ![samples/sonic.png](samples/sonic.png)

Results: ![samples/sonic_rotate.png](samples/sonic_rotate.gif)

### flash 

### tilemap
This option allows to create allow tiles and tile map from the input image. If for instance you use the mario level 1 like this : 
![samples/mario-level1.png](samples/mario-level1.png)
```martine -i mario-level1.png -w 16 -h 16 -tilemap -o Mario-level1```

By this command line, martine will analyse the image to find all tiles 16 pixels high and 16 pixels wide. It will keep in memory the tile position to produce the tile map.
![samples/mario-level1/tilesmap_schema.png](samples/mario-level1/tilesmap_schema.png)
A csv ot the same tile map is also created (tilesmap.map).
Martine will also generate all tiles found.
- ![samples/mario-level1/00.png](samples/mario-level1/00.png) tile 00 
- ![samples/mario-level1/01.png](samples/mario-level1/01.png) tile 01
- ![samples/mario-level1/02.png](samples/mario-level1/02.png) tile 02
- ![samples/mario-level1/03.png](samples/mario-level1/03.png) tile 03
- ![samples/mario-level1/04.png](samples/mario-level1/04.png) tile 04
- ![samples/mario-level1/05.png](samples/mario-level1/05.png) tile 05 
- ![samples/mario-level1/06.png](samples/mario-level1/06.png) tile 06
- ![samples/mario-level1/07.png](samples/mario-level1/07.png) tile 07
- ![samples/mario-level1/08.png](samples/mario-level1/08.png) tile 08
- ![samples/mario-level1/09.png](samples/mario-level1/09.png) tile 09
- ![samples/mario-level1/10.png](samples/mario-level1/10.png) tile 10 
- ![samples/mario-level1/01.png](samples/mario-level1/11.png) tile 11
- ![samples/mario-level1/02.png](samples/mario-level1/12.png) tile 12
- ![samples/mario-level1/03.png](samples/mario-level1/13.png) tile 13
- ![samples/mario-level1/04.png](samples/mario-level1/14.png) tile 14
- ![samples/mario-level1/05.png](samples/mario-level1/15.png) tile 15 
- ![samples/mario-level1/06.png](samples/mario-level1/16.png) tile 16
- ![samples/mario-level1/07.png](samples/mario-level1/17.png) tile 17
- ![samples/mario-level1/08.png](samples/mario-level1/18.png) tile 18
- ![samples/mario-level1/09.png](samples/mario-level1/19.png) tile 19
- ![samples/mario-level1/10.png](samples/mario-level1/20.png) tile 20 
- ![samples/mario-level1/01.png](samples/mario-level1/21.png) tile 21
- ![samples/mario-level1/02.png](samples/mario-level1/22.png) tile 22
- ![samples/mario-level1/03.png](samples/mario-level1/23.png) tile 23
- ![samples/mario-level1/04.png](samples/mario-level1/24.png) tile 24
- ![samples/mario-level1/05.png](samples/mario-level1/25.png) tile 25 
- ![samples/mario-level1/06.png](samples/mario-level1/26.png) tile 26
- ![samples/mario-level1/07.png](samples/mario-level1/27.png) tile 27
- ![samples/mario-level1/08.png](samples/mario-level1/28.png) tile 28
- ![samples/mario-level1/09.png](samples/mario-level1/29.png) tile 29
- ![samples/mario-level1/10.png](samples/mario-level1/30.png) tile 30 
- ![samples/mario-level1/10.png](samples/mario-level1/31.png) tile 31

and the palette for each tile for instance : ![samples/mario-level1/31_palettepal.png](samples/mario-level1/31_palettepal.png)
### Egx
The egx mode was introduced by Targhan in his game Ishido. 
This mode alternate differents screen mode. First line in mode 1, second line mode in mode 0, third in mode 1 etc ...

Two differents formats are available egx1 and egx2 : 
* egx1 alternates mode 0 and mode 1
* egx2 alternates mode 1 and mode 2

Martine lets the user to choose the input images. Like this you can combine dithering in the image mode 1, to get a better render.

In this example I will use this image : 
![bwind.jpg](samples/bwind.jpg)

So you need to generate two images in differents modes such as follow : 
```
martine -i images/bwind.jpg -m 0 -o egx_bwind/mode0 -f
martine -i images/bwind.jpg -m 1 -dithering 10 -o egx_bwind/mode1 -f
```
results : 
* mode 0 ![mode 0](samples/egx_bwind/mode0/bwind.jpg_down.png)
* mode 1 ![mode 1](samples/egx_bwind/mode1/bwind.jpg_down.png)

It's important to generate in differents folders, martine will erase the files generated by the first command.

Now create the egx file by the command : 
```
martine -egx1 -m 0 -i egx_bwind/mode0/BWINND.SCR -m2 1 egx_bwind/mode1/BWINND.SCR -o egx_bwind/ -pal egx_bwind/mode0/BWINND.PAL -f -dsk
```
will produce ![](samples/bwind_egx.png)

Some explanations : <br>
I choose here to be in full screen (option -f).
martine needs to know the two modes of each input images (here -m 0 mode 0 for the first image, -m2 1 mode 1 for the second image).<br>
-i and -i2 for the .scr images.<br>
And the palette path option -pal.<br>
You can note that martine allows you to choose you own palette. You can modify your palette to have a different rendering.<br>
If you want to iterate and gets the results quickly on your own machine, I advise you to add the M4 option and use the sna output. <br>
Complete your command line with for instance ```-sna -host 192.168.1.100 -autoexec ``` 

### deltapacking

This technic allows you to display animation on your CPC easily.<br>
For this technic you will need to script a little, even martine facilitates the data transformation.<br>
The steps are :
1. images traitment (if your starting image is gif for instance)
2. create a reference palette (this palette will be applied on all sprite conversion)
3. conversion of the sprite
4. differential computing on each sprite
5. parsing results 

For instance here, I will start with this gif image : 
![megaman](samples/deltapacking-megaman/megaman.gif)

You may not know, but the gif format allows pixel delta between sequence.
To avoid to get part of image, you need to traite the image before. 
I use here convert from ImageMagick : <br>


```convert megaman.gif -coalesce -scale 528x528 output.gif ```

You obtain this 
![megaman](samples/deltapacking-megaman/output.gif)

No visual changes, but now you can extract and get the whole image by sequence. 
<br>
Extract image from gif  (same convert binary): <br>

```convert output.gif images/m%02d.png```

Now create your reference palette :<br>

```martine -m 0 -i images/m00.png -o reference -dsk```

Iterate on all png to create sprite for each images : <br>

```for i in images/m*.png; do martine -i "$i" -m 0 -w 50 -h 50 -o sprites -pal reference/M000.PAL ; done ``` 

You will obtain in the sprites folder all sprites.

Compute the difference between images and ouput screen address (here the sprite will start at the address #D005): 

```martine -delta -df sprites/\*.WIN -o delta -address "#D005"```

Format your data : 

```prepare_delta -sprite sprites/m00.json -m 0 -delta delta/\*.json -out data.asm ```

The assembly code can be compiled with Rasm from Roudoudou : 
```
;--- dimensions du sprite ----
large equ 100 / 2 
haut equ 100
;-----------------------------

org #1000
run $

start

    ; gestion du mode 
    ;ld bc,#7f8c
    ;out (c),c 
;--- selection du mode ---------
    ld a,0
    call #BC0E
;-------------------------------
  

;--- gestion de la palette ---- 
    call palettefirmware
;------------------------------

call xvbl

;--- affichage du sprite initiale --  
    ; affichage du premier sprite
    ld de,#C000 ; adresse de l'ecran 
    ld hl,sprite ; pointeur sur l'image en memoire 
    ld b, haut ; hauteur de l'image 
    loop 
    push bc ; sauve le compteur hauteur dans la pile 
    push de ; sauvegarde de l'adresse ecran dans la pile
    ld bc, large ; largeur de l'image a afficher
    ldir ; remplissage de n * largeur octets a l'adresse dans de 
    pop de ; recuperation de l'adresse d'origine 
    ex de,hl ; echange des valeurs des adresses
    call bc26 ; calcul de l'adresse de la ligne suivante
    ex de,hl ; echange des valeurs des adresses
    pop bc ; retabli le compteur 
    djnz loop
;------------------------------------
   

mainloop    ; routine pour afficher les deltas provenant de martine 

;call #bb06
call xvbl
ld hl,delta00
call delta

;call #bb06

call xvbl
ld hl,delta01
call delta

;call #bb06
call xvbl
ld hl,delta02
call delta

;call #bb06
call xvbl
ld hl,delta03
call delta

;call #bb06
call xvbl
ld hl,delta04
call delta

;call #bb06
call xvbl
ld hl,delta05
call delta

;call #bb06
call xvbl
ld hl,delta06
call delta

;call #bb06
call xvbl
ld hl,delta07
call delta

;call #bb06
call xvbl
ld hl,delta08
call delta

;call #bb06
call xvbl
ld hl,delta09
call delta



jp mainloop 


;--- routine de deltapacking --------------------------

delta
 ld a,(hl) ; nombre de byte a poker
 ld (nbbytepoked),a ; stockage en mémoire
 inc hl
init
 ld a,(hl) ; octet a poker
 ld (pixel),a
 inc hl
 ld c,(hl) ; nbfois
 inc hl 
 ld b,(hl)
 inc hl
;
poke_octet
 ld e,(hl)
 inc hl
 ld d,(hl) ; de=adresse
 inc hl
 ld a,(pixel)
 ld (de),a ; poke a l'adresse dans de
 dec bc
 ld a,b ; test a t'on poke toutes les adresses compteur bc
 or a 
 jr nz, poke_octet
 ld a,c 
 or a
 jr nz, poke_octet

 ld a,(nbbytepoked) ; reste t'il d'autres bytes a poker ? 
 dec a 
 ld (nbbytepoked),a
 jr nz,init
 ret



;---------------------------------------------------------------
;
; attente de plusieurs vbl
;
xvbl ld e,30
	call waitvbl
	dec e
	jr nz,xvbl+2
	ret
;-----------------------------------

;---- attente vbl ----------
waitvbl
    ld b,#f5 ; attente vbl
vbl     
    in a,(c)
    rra
    jp nc,vbl
    ret
;---------------------------

;--- application palette firmware -------------
palettefirmware ; hl pointe sur les valeurs de la palette
ld a,0
ld b,0
ld c,0
call #bc32

ld a,1
ld b,15
ld c,15
call #bc32

ld a,2
ld b,1
ld c,1
call #bc32

ld a,3
ld b,10
ld c,10
call #bc32

ld a,4
ld b,11
ld c,11
call #bc32

ld a,5
ld b,12
ld c,12
call #bc32

ld a,6
ld b,13
ld c,13
call #bc32

ld a,7
ld b,14
ld c,14
call #bc32

ld a,8
ld b,23
ld c,23
call #bc32

ld a,9
ld b,16
ld c,16
call #bc32

ld a,10
ld b,26
ld c,26
call #bc32


ret
;---------------------------------------------

;---- recuperation de l'adresse de la ligne en dessous ------------
bc26 
ld a,h
add a,8 
ld h,a ; <---- le fameux que tu as oublié !
ret nc 
ld bc,#c050 ; on passe en 96 colonnes
add hl,bc
res 3,h
ret
;-----------------------------------------------------------------


;--- variables memoires -----
pixel db 0 
nbbytepoked db 0
;----------------------------


;------- data ---------------------------

include 'data.asm'



end

save'delta.bin',#1000,end-start,DSK,'bomberman.dsk'

```

Now you can get the result here : [sna](samples/deltapacking-megaman/megaman.sna) or [dsk](samples/deltapacking-megaman/megaman.dsk)

You will obtain this : 
![video emu](samples/deltapacking-megaman/megaman-emulator.gif)