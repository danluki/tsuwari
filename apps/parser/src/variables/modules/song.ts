import { PrismaService } from '@tsuwari/prisma';
import { UserIntegration } from '@tsuwari/spotify';

import { app } from '../../index.js';
import { LastFMDbData, LastFmIntegration } from '../../integrations/lastfm.js';
import { SpotifyIntegration } from '../../integrations/spotify.js';
import { VKDBData, VkIntegration } from '../../integrations/vk.js';
import { Module } from '../index.js';

const prisma = app.get(PrismaService);
const Lastfm = app.get(LastFmIntegration);
const VK = app.get(VkIntegration);
const Spotify = app.get(SpotifyIntegration);

export const song: Module = {
  key: 'currentsong',
  description: 'Current played song',
  handler: async (_, state) => {
    const enabledIntegrations = await prisma.channelIntegration.findMany({
      where: {
        channelId: state.channelId,
        enabled: true,
      },
      include: {
        integration: true,
      },
    });

    let result: string | null = '';

    for (const integration of enabledIntegrations) {
      const service = integration.integration.service;
      if (service === 'SPOTIFY') {
        if (!integration.accessToken || !integration.refreshToken) continue;
        const instance = await new UserIntegration(state.channelId, Spotify, prisma).init(integration);
        if (!instance) continue;
        const song = await instance.getCurrentSong();
        if (song) {
          result = song;
          break;
        } else continue;
      } else if (service === 'VK') {
        const data = integration.data as unknown as VKDBData;
        if (!data.userId) continue;
        const song = await VK.fetchSong(data.userId);
        if (song) {
          result = song;
          break;
        } else continue;
      } else if (service === 'LASTFM') {
        const data = integration.data as unknown as LastFMDbData;
        if (!data.username) continue;
        const song = await Lastfm.fetchSong(data.username);
        if (song) {
          result = song;
          break;
        } else continue;
      } else {
        result = null;
      }
    }

    if (result && result.length) {
      return result;
    } else {
      return;
    }
  },
};
