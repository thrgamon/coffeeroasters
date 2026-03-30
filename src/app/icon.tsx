import { ImageResponse } from 'next/og';

export const size = { width: 32, height: 32 };
export const contentType = 'image/png';

export default function Icon() {
	return new ImageResponse(
		<div
			style={{
				fontSize: 18,
				width: '100%',
				height: '100%',
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'center',
				background: '#280003',
				color: '#ffd500',
				fontWeight: 900,
				borderRadius: 4,
			}}
		>
			CR
		</div>,
		{ ...size },
	);
}
