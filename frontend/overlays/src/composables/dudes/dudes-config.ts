import { DudesSprite } from '@twir/types/overlays';
import type { SoundAsset, DudeAsset, AssetsLoadOptions } from '@twirapp/dudes/types';

export const dudesTwir = 'Twir';
export const dudesSprites = Object.keys(DudesSprite)
	.filter(sprite => sprite !== 'random') as (keyof typeof DudesSprite)[];

export function getRandomSprite() {
	return dudesSprites[Math.floor(Math.random() * dudesSprites.length)];
}

export const assetsLoadOptions: AssetsLoadOptions = {
  basePath: location.origin + '/overlays/dudes/sprites/',
  defaultSearchParams: {
    ts: Date.now(),
  },
};

export const dudesAssets: DudeAsset[] = [
  {
    alias: 'dude',
    src: 'dude/dude.json',
  },
  {
    alias: 'sith',
    src: 'sith/sith.json',
  },
  {
    alias: 'agent',
    src: 'agent/agent.json',
  },
  {
    alias: 'girl',
    src: 'girl/girl.json',
  },
  {
    alias: 'cat',
    src: 'cat/cat.json',
  },
  {
    alias: 'santa',
    src: 'santa/santa.json',
  },
];

export const dudesSounds: SoundAsset[] = [
	{
		alias: 'jump',
		src: location.origin + '/overlays/dudes/sounds/jump.mp3',
	},
];
