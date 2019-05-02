# martine

After Claudia from Eliot a long time ago, here is coming Martine.
A cli which converts JPEG or PNG image file into  SCR/PAL or Overscan  Amstrad CPC file screen.
Multi os, you can convert any pictures to Amstrad CPC Screen.
The files generated (.win, .scr, .ink) are compatible with [OCP art studio](http://www.cpc-power.com/index.php?page=detail&num=4963) and Impdraw V2 [i2](http://amstradplus.forumforever.com/t462-iMPdraw-v2-0.htm)

To Install and compile
```
go get github.com/jeromelesaux/martine
cd $GOPATH/src/github.com/jeromelesaux/martine
go get 
go build
```

To get binary : 
[https//github.com/jeromelesaux/martine/releases](https//github.com/jeromelesaux/martine/releases)
<br>OS avaible : Linux, Macos X and Windows  

Usage and options : 

```
./martine
martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)
By Impact Sid (Version:0.16.rc)
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
  -autoexec
    	Execute on your remote CPC the screen file or basic file.
  -dsk
    	Copy files in a new CPC image Dsk.
  -f	Overscan mode (default no overscan)
  -h int
    	Custom output height in pixels. (default -1)
  -help
    	Display help message
  -host string
    	Set the ip of your M4.
  -i string
    	Picture path of the input file.
  -info
    	Return the information of the file, associated with -pal and -win options
  -ink string
    	Path of the palette Cpc ink file.
  -iter int
    	Iterations number to walk in roll mode, or number of images to generate in rotation mode. (default -1)
  -iterx int
    	Number of tiles on a row in the input image. (default -1)
  -itery int
    	Number of tiles on a column in the input image. (default -1)
  -keephigh int
    	bit rotation on the top and keep pixels (default -1)
  -keeplow int
    	bit rotation on the bottom and keep pixels (default -1)
  -kit string
    	Path of the palette Cpc plus Kit file.
  -losthigh int
    	bit rotation on the top and lost pixels (default -1)
  -lostlow int
    	bit rotation on the bottom and lost pixels (default -1)
  -m int
    	Output mode to use :
    		0 for mode0
    		1 for mode1
    		2 for mode2
    		and add -f option for overscan export.
    		 (default -1)
  -n	no amsdos header for all files (default amsdos header added).
  -o string
    	Output directory
  -p	Plus mode (means generate an image for CPC Plus Screen)
  -pal string
    	Apply the input palette to the image
  -remotepath string
    	remote path on your M4 where you want to copy your files.
  -rla int
    	bit rotation on the left and keep pixels (default -1)
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
    	bit rotation on the right and keep pixels (default -1)
  -s string
    	Byte statement to replace in ascii export (default is BYTE), you can replace or instance by defb
  -sla int
    	bit rotation on the left and lost pixels (default -1)
  -sra int
    	bit rotation on the right and lost pixels (default -1)
  -tile
    	Tile mode to create multiples sprites from a same image.
  -w int
    	Custom output width in pixels. (default -1)
  -win string
    	Filepath of the ocp win file
  -z int
    	Compression algorithm : 
    		1: rle (default)
    		2: rle 16bits
    		3: Lz4 Classic
    		4: Lz4 Raw
    	 (default -1)

```

Principles : 
Martine can be used in 5 modes : 
 * first mode : conversion of an input image into sprite (file .win), cpc screen (file .scr) and overscan screen (larger file .scr)
 * second mode : rotate line or column pixels of an input image (option -roll)
 * bulk conversion of a tiles page : each sprites of the same width and height will be converted to sprites (files .win, option -tile)
 * rotate an existing image, will produce images rotated
 * 3d rotate an existing image on an axis.

examples :
## 1. Screen conversion :
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

 files generated : 
 * .win or .scr sprite or screen files
 * .pal or .ink palette file (.ink will be generated if the -p option is set)
 * .txt ascii file with palettes values (firmware values and basic values), and screen byte values
 *  c.txt ascii file with palettes values (firmware values and basic values), and screen byte values by column
 * .json json file with palettes values (firmware values and basic values), and screen byte values
 * .bas launch to test the screen load on classic .scr 17ko 
 * _resized.png images files to ensure the resize action
 * _downgraded.png images files to ensure the downgraded palette action

additionnals options available : 
* -dsk will generate a dsk file and add all amsdos files will be added.
* -n will remove amsdos headers from the amsdos files
* -f will generate overscan screen amsdos file
* -p will generate a CPC plus screen amsdos file
* -h will generate sprite of x pixel high
* -w will generate sprite of x pixel wide
* -s to define the byte token will be replace the byte token in the ascii files
* -a to set the algorithm to downsize the image
* -m to define the screen mode 0,1,2
* -o to set the output directory
* -z compress the image .scr / .win (not overscan) using the algorithm

## 2. Samples roll usage : 

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

 files generated : 
 * .win sprite files
 * .pal or .ink palette file (.ink will be generated if the -p option is set)
 * .txt ascii file with palettes values (firmware values and basic values), and screen byte values
 *  c.txt ascii file with palettes values (firmware values and basic values), and screen byte values by column
 * .json json file with palettes values (firmware values and basic values), and screen byte values
 *  _resized.png images files to ensure the resize action
 * _downgraded.png images files to ensure the downgraded palette action

additionnals options available : 
* -dsk will generate a dsk file and add all amsdos files will be added.
* -n will remove amsdos headers from the amsdos files
* -f will generate overscan screen amsdos file
* -p will generate a CPC plus screen amsdos file
* -h will generate sprite of x pixel high
* -w will generate sprite of x pixel wide
* -m to define the screen mode 0,1,2
* -o to set the output directory
* -sla will rotate x column pixels from the left, those columns will be discarded  
* -sra will rotate x column pixels from the right, those columns will be discarded  
* -rra will rotate x column pixels from the right
* -rla will rotate x column pixels from the left
* -keephigh will rotate x line pixels to the top
* -keeplow will rotate x line pixels to the bottom
* -losthigh will rotate x line pixels to the top, those lines will be discarded 
* -lostlow will rotate x line pixels to the bottom, those lines will be discarded
* -s to define the byte token will be replace the byte token in the ascii files
* -a to set the algorithm to downsize the image
* -z compress the image .scr / .win (not overscan) using the algorithm

## 3. tile option :

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

  files generated : 
 * .win sprite files
 * .pal or .ink palette file (.ink will be generated if the -p option is set)
 * .txt ascii file with palettes values (firmware values and basic values), and screen byte values
 *  c.txt ascii file with palettes values (firmware values and basic values), and screen byte values by column
 * .json json file with palettes values (firmware values and basic values), and screen byte values
 *  _resized.png images files to ensure the resize action
 * _downgraded.png images files to ensure the downgraded palette action

additionnals options available : 
* -dsk will generate a dsk file and add all amsdos files will be added.
* -n will remove amsdos headers from the amsdos files
* -f will generate overscan screen amsdos file
* -p will generate a CPC plus screen amsdos file
* -h will generate sprite of x pixel high
* -w will generate sprite of x pixel wide
* -m to define the screen mode 0,1,2
* -o to set the output directory
* -s to define the byte token will be replace the byte token in the ascii files
* -a to set the algorithm to downsize the image
* -z compress the image .scr / .win (not overscan) using the algorithm

## 4. rotation sprite : 
This option able to rotate the input image, iter will generate the number of images.
```martine -i images/coke.jpg -rotate -iter 16 -o test  -m 1 -w 32 -h 32```

Input : ![samples/coke.jpg](samples/coke.jpg)

results : ![samples/coke.igf](samples/coke.gif)

  files generated : 
 * .win sprite files
 * .pal or .ink palette file (.ink will be generated if the -p option is set)
 * .txt ascii file with palettes values (firmware values and basic values), and screen byte values
 *  c.txt ascii file with palettes values (firmware values and basic values), and screen byte values by column
 * .json json file with palettes values (firmware values and basic values), and screen byte values
 *  _resized.png images files to ensure the resize action
 * _downgraded.png images files to ensure the downgraded palette action

additionnals options available : 
* -dsk will generate a dsk file and add all amsdos files will be added.
* -n will remove amsdos headers from the amsdos files
* -f will generate overscan screen amsdos file
* -p will generate a CPC plus screen amsdos file
* -h will generate sprite of x pixel high
* -w will generate sprite of x pixel wide
* -m to define the screen mode 0,1,2
* -o to set the output directory
* -s to define the byte token will be replace the byte token in the ascii files
* -a to set the algorithm to downsize the image
* -z compress the image .scr / .win (not overscan) using the algorithm

## 5. 3d rotation sprite : 
This option able to rotate the input image on an axis (X or Y and even on a diagonal), iter will generate the number of images. 
```./martine -m 0 -w 64 -h 64 -rotate3d -rotate3dtype 2 -o test/ -i images/sonic.png  -iter 12```

Input : ![samples/sonic.png](samples/sonic.png)

Results: ![samples/sonic_rotate.png](samples/sonic_rotate.gif)

 files generated : 
 * .win sprite files
 * .pal or .ink palette file (.ink will be generated if the -p option is set)
 * .txt ascii file with palettes values (firmware values and basic values), and screen byte values
 *  c.txt ascii file with palettes values (firmware values and basic values), and screen byte values by column
 * .json json file with palettes values (firmware values and basic values), and screen byte values
 *  _resized.png images files to ensure the resize action
 * _downgraded.png images files to ensure the downgraded palette action

additionnals options available : 
* -dsk will generate a dsk file and add all amsdos files will be added.
* -n will remove amsdos headers from the amsdos files
* -f will generate overscan screen amsdos file
* -p will generate a CPC plus screen amsdos file
* -h will generate sprite of x pixel high
* -w will generate sprite of x pixel wide
* -m to define the screen mode 0,1,2
* -o to set the output directory
* -s to define the byte token will be replace the byte token in the ascii files
* -a to set the algorithm to downsize the image
* -z compress the image .scr / .win (not overscan) using the algorithm