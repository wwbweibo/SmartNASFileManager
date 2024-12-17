# SmartNASFileManager

This project aims using AI to simplify the file management on NAS.

### extra things when use this project

decord install failed when install LAVIS:

decord provided pypi package not supported for arm. so need to build it yourself. [decord](https://github.com/dmlc/decord).
 
While build decord yourself, may meet build failed caused by ffmpeg version. Make sure that you have ffmpeg 4.X version installable, after this, rerun make with flag ` -DFFMPEG_DIR='[your ffmpeg 4.X install path]'`