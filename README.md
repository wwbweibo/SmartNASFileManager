# SmartNASFileManager

This project aims using AI to simplify the file management on NAS.

## TODO List

- File System Base:
  - [x] File Browser
  - [ ] File Manual Tag 
  - [ ] FileIndex create and update at realtime 
  - [ ] FileSystem event watch
  - [ ] File encryption at write
  - [ ] Multi NAS Support
- Image Files:
  - [X] Image Browser
  - [X] Image Snapshot
  - [X] Image Caption Using Local Vision Model
  - [X] Image Auto Tag
  - [ ] Image caption and tag using LLM
  - [ ] Image search by tag and caption
  - [ ] Image search by similar
  - [ ] RAW Image Support
- Video Files: 
  - [ ] Video Player
  - [ ] Video Caption
- Document Files:
  - [ ] Support Document preview and edit
  - [ ] Using RAG to build knowledge base
  - [ ] Document search using vec

### extra things when use this project

decord install failed when install LAVIS:

decord provided pypi package not supported for arm. so need to build it yourself. [decord](https://github.com/dmlc/decord).
 
While build decord yourself, may meet build failed caused by ffmpeg version. Make sure that you have ffmpeg 4.X version installable, after this, rerun make with flag ` -DFFMPEG_DIR='[your ffmpeg 4.X install path]'`