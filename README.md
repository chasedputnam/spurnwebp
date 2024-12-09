# spurnwebp
A golang based utility that can convert WEBP image files into PNG or JPG from the command line. The utility can also watch a directory and convert all WEBP images added into it into PNG or JPG.

Usage:

```aiignore
spurnwebp -input /home/username/images/example_a.webp 
```
**-input** will take in an absolute file path to a .webp image.

The utility will output a PNG file by default in the same directory as the input webp file.

```aiignore
spurnwebp -watch /home/username/images
```
**-watch** will create a filesystem watcher to the provided directory path.

Any files added to the directory ending in .webp will be detected and converted to PNG by default.

```aiignore
spurnwebp -input /home/username/images/example_a.webp -outputType JPG
```
**-input** will take in an absolute file path to a .webp image.

**-outputType** will accept PNG or JPG as values and allow the utility to convert to PNG or JPG formats.

The utility will output a JPG file in the same directory as the input webp file.

```aiignore
spurnwebp -watch /home/username/images - outputType JPG
```
**-watch** will create a filesystem watcher to the provided directory path.

**-outputType** will accept PNG or JPG as values and allow the utility to convert to PNG or JPG formats.


Any files added to the directory ending in .webp will be detected and converted to JPG.
