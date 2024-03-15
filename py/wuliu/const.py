Project_JSON = "project.json"
Thumb_Size    = "ThumbSize"

ID          = 'ID'
Filename    = 'Filename'
Checksum    = 'Checksum'
Size        = 'Size'
Type        = 'Type'
Like        = 'Like'
Label       = 'Label'
Notes       = 'Notes'
Keywords    = 'Keywords'
Collections = 'Collections'
Albums      = 'Albums'
CTime       = 'CTime'
UTime       = 'UTime'


Files       = 'files'
Thumbs      = 'webpages/thumbs'
Thumbs_msgp = 'thumbs.msgp'


New_Album_Info = {  # 用於創建新相簿
    ids: [],        # 通過 ID 指定圖片 (一旦指定ID, 其他條件無效)
    label: '',      # label, notes, keyword, collections, albums
    notes: '',      # 這五項可取交集(默認), 也可取併集, 通過下面的 union 控制。
    keywords: [],
    collections: [],
    albums: [],
    union: False,  # False: 交集, True: 併集(聯集)
    orderby: 'utime',  # 排序依據: utime/ctime/filename/like
    ascending: False,  # False: 降序, True: 昇序, 如果指定 ids, 則以 ids 為準
}

