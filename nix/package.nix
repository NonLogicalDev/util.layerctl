{ lib
, buildGoModule
, git
, makeWrapper
, source ? lib.fileset.toSource {
    root = ../.;
    fileset = lib.fileset.unions [ ../cmd ../internal ../go.mod ../go.sum ];
  }
, revision ? "unknown"
}:

buildGoModule (finalAttrs: {
  pname = "layerctl";
  version = "0.2.0";
  src = source;
  vendorHash = "sha256-OK7q/CQpa4gNWIPWHZPyew2Ym8NPkXqSDSzE92gYrKk=";
  ldflags = [ "-s" "-w" "-X main.version=${finalAttrs.version}" "-X main.commit=${revision}" ];
  subPackages = [ "cmd/layerctl" ];
  nativeBuildInputs = [ makeWrapper ];
  nativeCheckInputs = [ git ];

  postInstall = ''
    wrapProgram "$out/bin/layerctl" --suffix PATH : ${lib.makeBinPath [ git ]}
  '';

  checkPhase = ''
    runHook preCheck
    go test ./...
    runHook postCheck
  '';

  meta = {
    description = "Compose Git repositories and changes as ordered layers";
    homepage = "https://github.com/NonLogicalDev/util.layerctl";
    license = lib.licenses.asl20;
    mainProgram = "layerctl";
  };
})
