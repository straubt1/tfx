// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://tfx.rocks',
  integrations: [
    starlight({
      title: 'TFx',
      description: 'A standalone CLI for HCP Terraform and Terraform Enterprise.',
      favicon: '/favicon.ico',
      social: {
        github: 'https://github.com/straubt1/tfx',
      },
      customCss: ['./src/styles/custom.css'],
      sidebar: [
        {
          label: 'Commands',
          items: [
            { label: 'Organization', slug: 'commands/organization' },
            { label: 'Project', slug: 'commands/project' },
            { label: 'Variable Sets', slug: 'commands/variable_set' },
            {
              label: 'Workspace',
              items: [
                { label: 'General', slug: 'commands/workspace' },
                { label: 'Plans', slug: 'commands/workspace_plan' },
                { label: 'Runs', slug: 'commands/workspace_run' },
                { label: 'Variables', slug: 'commands/workspace_variable' },
                { label: 'Configuration Versions', slug: 'commands/workspace_configurationversion' },
                { label: 'Team', slug: 'commands/workspace_team' },
                { label: 'Lock', slug: 'commands/workspace_lock' },
                { label: 'State Versions', slug: 'commands/workspace_stateversion' },
              ],
            },
            {
              label: 'Private Registry',
              items: [
                { label: 'Modules', slug: 'commands/registry_module' },
                { label: 'Providers', slug: 'commands/registry_provider' },
              ],
            },
            { label: 'Releases', slug: 'commands/release' },
            {
              label: 'Admin',
              items: [
                { label: 'GPG Keys', slug: 'commands/admin_gpg' },
                { label: 'Terraform Versions', slug: 'commands/admin_terraformversion' },
              ],
            },
          ],
        },
        {
          label: 'Integrations',
          items: [{ label: 'GitLab CI', slug: 'integrations/gitlab' }],
        },
        {
          label: 'Debugging',
          items: [
            { label: 'TFX_LOG', slug: 'debugging/log-level' },
            { label: 'TFX_LOG_PATH', slug: 'debugging/log-path' },
          ],
        },
        {
          label: 'Testing',
          items: [{ label: 'Test Plan', slug: 'testing/test_plan' }],
        },
        {
          label: 'About',
          items: [
            { label: 'Why TFx?', slug: 'about/purpose' },
            {
              label: 'Release Notes',
              link: 'https://github.com/straubt1/tfx/blob/main/CHANGELOG.md',
              attrs: { target: '_blank', rel: 'noopener noreferrer' },
            },
            { label: 'Contributing', slug: 'about/contributing' },
            { label: 'Style Guide', slug: 'about/style_guide' },
          ],
        },
      ],
    }),
  ],
});
