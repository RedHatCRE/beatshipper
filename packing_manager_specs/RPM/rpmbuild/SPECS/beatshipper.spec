%define name    beatshipper
%define version 0.0.1
%define _rpmdir %{getenv:GITHUB_WORKSPACE}/

Name:           %{name}
Version:        %{version}
Release:        1%{?dist}
Summary:        GNU ZIP beats shipper
License:        GPL
Requires:       bash
URL:            https://github.com/RedHatCRE/%{name}/

%description
Since there’s no way to send GNU zip files through filebeat, this service will be responsible for checking if there are new .gz files based on a path that we’ll explode using globbing, decompress them, send them using the filebeat service with the provided configuration and store them in a local file registry.

%prep

%install

# Create folders of binary and config dir with right permissions:
install -m 0755 -d %{buildroot}/etc/%{name}
install -m 0755 -d %{buildroot}/usr/sbin/
install -m 0755 -d %{buildroot}/lib/systemd/system/

# Copy files with right permissions:
install -m 0755 %{_rpmdir}%{name} %{buildroot}/usr/sbin/%{name}
install -m 0644 %{_rpmdir}/beatshipper-conf.yml %{buildroot}/etc/%{name}/
install -m 0644 %{_rpmdir}/lib/systemd/system/beatshipper.service %{buildroot}/lib/systemd/system/

%post

%clean
rm -rf $RPM_BUILD_ROOT

%files
/usr/sbin/%{name}
/etc/%{name}/beatshipper-conf.yml
/lib/systemd/system/beatshipper.service

%changelog
 # END SPEC