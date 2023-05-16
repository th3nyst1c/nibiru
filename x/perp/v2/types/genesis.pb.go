// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: perp/v2/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the perp module's genesis state.
type GenesisState struct {
	Markets          []Market          `protobuf:"bytes,2,rep,name=markets,proto3" json:"markets"`
	Amms             []AMM             `protobuf:"bytes,3,rep,name=amms,proto3" json:"amms"`
	Positions        []Position        `protobuf:"bytes,4,rep,name=positions,proto3" json:"positions"`
	ReserveSnapshots []ReserveSnapshot `protobuf:"bytes,5,rep,name=reserve_snapshots,json=reserveSnapshots,proto3" json:"reserve_snapshots"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8edcabc35f3cf683, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetMarkets() []Market {
	if m != nil {
		return m.Markets
	}
	return nil
}

func (m *GenesisState) GetAmms() []AMM {
	if m != nil {
		return m.Amms
	}
	return nil
}

func (m *GenesisState) GetPositions() []Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *GenesisState) GetReserveSnapshots() []ReserveSnapshot {
	if m != nil {
		return m.ReserveSnapshots
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "nibiru.perp.v2.GenesisState")
}

func init() { proto.RegisterFile("perp/v2/genesis.proto", fileDescriptor_8edcabc35f3cf683) }

var fileDescriptor_8edcabc35f3cf683 = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x4f, 0x4b, 0x03, 0x31,
	0x10, 0xc5, 0xb7, 0x7f, 0x54, 0x5c, 0x45, 0x74, 0xab, 0xb2, 0x14, 0x49, 0xc5, 0x93, 0x17, 0x37,
	0xb4, 0x05, 0x2f, 0x7a, 0xb1, 0x1e, 0x3c, 0xad, 0x48, 0x7b, 0xf3, 0x22, 0xd9, 0x12, 0xd2, 0xa0,
	0x9b, 0x09, 0x99, 0x74, 0xd1, 0x4f, 0xe0, 0xd5, 0x8f, 0xd5, 0x63, 0x8f, 0x9e, 0x44, 0xda, 0x2f,
	0x22, 0xcd, 0xa6, 0x48, 0xf7, 0x12, 0xc2, 0xbc, 0xf7, 0x7b, 0xf3, 0x98, 0xf0, 0x44, 0x73, 0xa3,
	0x69, 0xd1, 0xa3, 0x82, 0x2b, 0x8e, 0x12, 0x13, 0x6d, 0xc0, 0x42, 0x74, 0xa0, 0x64, 0x26, 0xcd,
	0x34, 0x59, 0xa9, 0x49, 0xd1, 0x6b, 0x1f, 0x0b, 0x10, 0xe0, 0x24, 0xba, 0xfa, 0x95, 0xae, 0xf6,
	0x99, 0x00, 0x10, 0x6f, 0x9c, 0x32, 0x2d, 0x29, 0x53, 0x0a, 0x2c, 0xb3, 0x12, 0x94, 0xcf, 0x68,
	0x93, 0x31, 0x60, 0x0e, 0x48, 0x33, 0x86, 0x9c, 0x16, 0xdd, 0x8c, 0x5b, 0xd6, 0xa5, 0x63, 0x90,
	0xca, 0xeb, 0xad, 0xf5, 0x6a, 0xb4, 0xcc, 0xf2, 0x72, 0x78, 0xf1, 0x59, 0x0f, 0xf7, 0x1f, 0xca,
	0x2a, 0xa3, 0xd5, 0x38, 0xba, 0x0e, 0x77, 0x72, 0x66, 0x5e, 0xb9, 0xc5, 0xb8, 0x7e, 0xde, 0xb8,
	0xdc, 0xeb, 0x9d, 0x26, 0x9b, 0xdd, 0x92, 0xd4, 0xc9, 0x83, 0xe6, 0xec, 0xa7, 0x13, 0x0c, 0xd7,
	0xe6, 0xe8, 0x2a, 0x6c, 0xb2, 0x3c, 0xc7, 0xb8, 0xe1, 0xa0, 0x56, 0x15, 0xba, 0x4b, 0x53, 0x4f,
	0x38, 0x5b, 0x74, 0x1b, 0xee, 0x6a, 0x40, 0xe9, 0xfa, 0xc7, 0x4d, 0xc7, 0xc4, 0x55, 0xe6, 0xc9,
	0x1b, 0x3c, 0xf8, 0x0f, 0x44, 0xc3, 0xf0, 0xc8, 0x70, 0xe4, 0xa6, 0xe0, 0x2f, 0xa8, 0x98, 0xc6,
	0x09, 0x58, 0x8c, 0xb7, 0x5c, 0x4a, 0xa7, 0x9a, 0x32, 0x2c, 0x8d, 0x23, 0xef, 0xf3, 0x61, 0x87,
	0x66, 0x73, 0x8c, 0x83, 0x74, 0xb6, 0x20, 0xb5, 0xf9, 0x82, 0xd4, 0x7e, 0x17, 0xa4, 0xf6, 0xb5,
	0x24, 0xc1, 0x7c, 0x49, 0x82, 0xef, 0x25, 0x09, 0x9e, 0xfb, 0x42, 0xda, 0xc9, 0x34, 0x4b, 0xc6,
	0x90, 0xd3, 0x47, 0x17, 0x7e, 0x3f, 0x61, 0x52, 0xd1, 0x72, 0x11, 0x7d, 0xa7, 0xeb, 0xc3, 0xda,
	0x0f, 0xcd, 0xf1, 0xc6, 0xbd, 0xd9, 0xb6, 0xbb, 0x6f, 0xff, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x07,
	0x8a, 0x78, 0xdd, 0xf1, 0x01, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ReserveSnapshots) > 0 {
		for iNdEx := len(m.ReserveSnapshots) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ReserveSnapshots[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Amms) > 0 {
		for iNdEx := len(m.Amms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Markets) > 0 {
		for iNdEx := len(m.Markets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Markets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Markets) > 0 {
		for _, e := range m.Markets {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Amms) > 0 {
		for _, e := range m.Amms {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ReserveSnapshots) > 0 {
		for _, e := range m.ReserveSnapshots {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Markets", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Markets = append(m.Markets, Market{})
			if err := m.Markets[len(m.Markets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amms", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amms = append(m.Amms, AMM{})
			if err := m.Amms[len(m.Amms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Positions = append(m.Positions, Position{})
			if err := m.Positions[len(m.Positions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReserveSnapshots", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReserveSnapshots = append(m.ReserveSnapshots, ReserveSnapshot{})
			if err := m.ReserveSnapshots[len(m.ReserveSnapshots)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
