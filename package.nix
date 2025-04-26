{buildGoModule}:
buildGoModule {
  pname = "gql";
  version = "0.1";

  src = ./.;

  checkPhase = ''
    go test -v -short ./...
  '';

  vendorHash = "sha256-5zShQBdpZ9RacQn1yKsxPAJMgfYKMM4WvRdVtRLLC/0=";
}
