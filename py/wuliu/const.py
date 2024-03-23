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
    name: ''        # 相簿名稱，必填，只允許使用 -, _, 0-9, a-z, A-Z
    ids: [],        # 通過 ID 指定圖片 (一旦指定ID, 其他條件無效)
    label: '',      # label, notes, keyword, collections, albums
    notes: '',      # 這五項可取併集(默認), 也可取交集, 通過下面的 union 控制。
    keywords: [],   # 這五項其中留空的，則被忽略。
    collections: [],  # 如果這五項及 ids 都留空，則表示 "全部圖片"。
    albums: [],
    union: True,       # True: 併集(聯集), False: 交集
    orderby: 'utime',  # 排序依據: utime/ctime/filename/like
    ascending: False,  # False: 降序, True: 昇序, 如果指定 ids, 則以 ids 為準
}

