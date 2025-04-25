{buildGoModule}:
buildGoModule {
  pname = "gql";
  version = "0.1";

  src = ./.;

  checkPhase = ''
    go test -v -short ./...
  '';

  vendorHash = "sha256-yNBzWcSZzb+yywLaOtBd9sS0DP8XfkUAsnpBykFp5jI=";
}
