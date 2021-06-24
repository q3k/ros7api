with import <nixpkgs> {};

pkgs.mkShell {
  buildInputs = with pkgs; [
    protobuf
    go
    go-protobuf
  ];
}
