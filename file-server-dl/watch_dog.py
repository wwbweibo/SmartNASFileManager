import sys
from watchdog.observers import Observer
from watchdog.events import  FileSystemEventHandler
import os
from os import path 
import time
import hashlib

def md5checksum(f: str):
    m = hashlib.md5()
    with open(f, 'rb') as file:
        while chunk := file.read(8192):
            m.update(chunk)
    return m.hexdigest()

class Handler(FileSystemEventHandler):
    def __init__(self):
        self.device_to_dir = {
            'X-T50': '/share/Public/images/x-t50',
            'GFX50S II': '/share/Public/images/gfx50s-ii'
        }
    def on_any_event(self, event):
        if event.event_type == 'created':
            print(f"event type: {event.event_type}, path: {event.src_path}")
        # if file created and name contains X-50, fuji x-t50 mounted, do file copy
        if event.event_type == "created" and "X-T50" in event.src_path:
            print(f"camera X-T50 mounted at {event.src_path}")
            self.copy_file("X-T50", event.src_path)
        if event.event_type == "created" and "GFX50S II" in event.src_path:
            print(f"camera GFX50S II mounted at {event.src_path}")
            self.copy_file("GFX50S II", event.src_path)
    def copy_file(self, device_name: str,  mounted_path: str):
        target_dir = self.device_to_dir[device_name]
        # reading all files from camera
        entries =  os.walk(mounted_path)
        print(f"ready to backup files to nas dir {target_dir}")
        for dir, subdir, files in entries:
            if len(files) > 0 and (files[0].endswith('.JPG') or files[0].endswith('RAF')):
                # exist files in camera
                # check and copy files
                for f in files:
                    src_path = path.join(dir, f)
                    # get file create time, and format into yyyy-mm-dd
                    create_time = path.getmtime(src_path)
                    create_date = time.strftime('%Y-%m-%d', time.localtime(create_time))
                    _target_dir = path.join(target_dir, create_date)
                    if not path.exists(_target_dir):
                        os.makedirs(_target_dir)
                    target_path = path.join(_target_dir, f)
                    # check if file already exist, if exists check hash sum is equal
                    if path.exists(target_path) and  md5checksum(src_path) == md5checksum(target_path):
                        print(f'file {src_path} already exists in {target_path}')
                        continue
                    os.system(f"cp '{src_path}' '{target_path}'")
                    print(f'copy {src_path} to {target_path}')
        print(f'copy files from {mounted_path} to {target_dir} done')

if __name__ == "__main__":
    watch_path = sys.argv[1] if len(sys.argv) > 1 else '.'
    event_handler = Handler()
    observer = Observer()
    observer.schedule(event_handler, watch_path, recursive=True)
    observer.start()
    try:
        while True:
            observer.join(1)
    finally:
        observer.stop()
        observer.join()