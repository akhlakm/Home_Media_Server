Run the media indexer directly on server to crawl the /media/i files
and move them to img/ or mp4/ or gif/ directories automatically.

A JS and JSON file will be created that can be used in the InfiniteScroller.html
to view the files over http.

Specify -serve argument to spawn a webserver on port 9000. Reverse proxy your nginx
server configuration (/etc/nginx/nginx.conf) on this port. Default home page is the 
infinitescroller page.

Specify -walk argument to actually start crawling the directories. Alternatively
get request /api/walk/ to start walking the directories from the frontend.

The application depends on ffprobe and ffmpeg to be available on the $PATH.
Make sure they are accessible in the terminal first.

During crawling, any non MP4 videos will be converted. Videos longer than 20 secs
will be segemented into MP4 videos. These ffmpeg processes will run in the background
that can take days to finish depending on the file types. The original files will be
removed upon successful conversion or segmentation.

In the foreground, the available images, gifs, and mp4 files will be moved and added
to the list.js and list.json files. Once ffmpeg finishes with conversion and segmentation,
similar listing will be rerun.

To build, run the build.ps1 powershell script, preferably in VSCode.

A new simpler JsCanvas is added to the frontend, which can be used to add caption to an
image directly from the infinite scroller. Make sure your server configuration accepts
file upload. (Increase nginx upload size limit as client_max_body_size 100M;).

