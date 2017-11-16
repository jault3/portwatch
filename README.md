# portwatch

> NOTE: this only works on macOS

Simple program that polls for open ports (binding to anything except localhost) and displays a notification when found.

## launchctl config

Recommended way to run is using launchctl:

```
go get -u github.com/jault3/portwatch

cat <<EOF > ~/Library/LaunchAgents/com.github.jault3.portwatch.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>portwatch</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{ GOBIN }}/portwatch</string>
    </array>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

launchctl load ~/Library/LaunchAgents/com.github.jault3.portwatch.plist
```

Be sure to replace {{ GOBIN }} with your absolute path to $GOBIN. You cannot use environment variables, relative paths, or binary names that depend on $PATH in the .plist file (unless of course you know how launchctl works and you set it up correctly).
