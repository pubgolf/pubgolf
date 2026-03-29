export function cleanPhoneNumber(num: string): string {
	num = num.replaceAll(/[^\d]/g, '');
	if (num === '') {
		return '';
	}
	if (num.length < 11 && !num.startsWith('1')) {
		num = '1' + num;
	}
	return '+' + num;
}
