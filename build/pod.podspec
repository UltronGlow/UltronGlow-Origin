Pod::Spec.new do |spec|
  spec.name         = 'utg'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/UltronGlow/UltronGlow-Origin'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS utg Client'
  spec.source       = { :git => 'https://github.com/UltronGlow/UltronGlow-Origin.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/utg.framework'

	spec.prepare_command = <<-CMD
    curl https://utgstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/utg.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
