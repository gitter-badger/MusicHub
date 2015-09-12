<html>
  <head>
    <title>Upload File</title>
  </head>
  <body>
    <form enctype="multipart/form-data" action="/upload" name="uploadsong" method="post">
      SongFile:<input type="file" name="upload_file"/>
      <input type="hidden" name="token" value="{{.}}"/>
      <input type="submit" value="Upload"/>
    </form>
  </body>
</html>
