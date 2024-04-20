import type { Message, PlainMessage } from '@bufbuild/protobuf';

export type Strict<T extends Message<T>> = Required<PlainMessage<T>>;
