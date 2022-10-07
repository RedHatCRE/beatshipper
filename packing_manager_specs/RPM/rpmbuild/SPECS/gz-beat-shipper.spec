%define name    gz-beat-shipper
%define version 0.0.1
%define _rpmdir %{getenv:GITHUB_WORKSPACE}/

Name:           %{name}
Version:        %{version}
Release:        1%{?dist}
Summary:        GNU ZIP beats shipper
License:        GPL
Requires:       bash
Source:         %{name}‑%{version}.tar.gz
URL:            https://github.com/RedHatCRE/gz-beat-shipper/

%description
Since there’s no way to send GNU zip files through filebeat, this service will be responsible for checking if there are new .gz files based on a path that we’ll explode using globbing, decompress them, send them using the filebeat service with the provided configuration and store them in a local file registry.

%prep

%install
echo ${GITHUB_WORKSPACE}

%post

%clean
rm -rf $RPM_BUILD_ROOT

%files

%changelog
 # END SPEC
