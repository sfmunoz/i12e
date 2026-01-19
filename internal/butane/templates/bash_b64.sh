base64 -d <<< "{{ .Buf }}" | gunzip | bash
