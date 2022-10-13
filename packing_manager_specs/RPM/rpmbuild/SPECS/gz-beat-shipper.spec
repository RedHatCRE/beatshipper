%define name    gz-beat-shipper
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

install -m 0755 -d %{buildroot}/etc/%{name}
install -m 0755 -d %{buildroot}/usr/sbin/
install -m 0644 %{_rpmdir}%{name} %{buildroot}/usr/sbin/%{name}
install -m 0644 %{_rpmdir}/gz-beat-shipper-conf.yml %{buildroot}/etc/%{name}/
install -m 0644 %{_rpmdir}/lib/systemd/system/gz-beat-shipper.service %{buildroot}/lib/systemd/system/

%post

%clean
rm -rf $RPM_BUILD_ROOT

%files
/usr/sbin/%{name}
/etc/%{name}/gz-beat-shipper-conf.yml
/lib/systemd/system/gz-beat-shipper.service

%changelog
 # END SPEC